package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"golang.org/x/sync/errgroup"
	"homework/cmd/cli"
	"homework/cmd/httpserv"
	"homework/internal/db"
	"homework/internal/logger"
	"homework/internal/model"
	"homework/internal/service/order"
	"homework/internal/service/pickuppoint"
	"homework/internal/storage/file"
	"homework/internal/storage/postgres"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"text/tabwriter"
	"time"
)

const (
	ORDERS_FILEPATH = "orders.json"
	POINTS_FILEPATH = "points.json"
	HTTPS_ADDR      = ":9443"
	REDIRECT_ADDR   = ":9000"
	CERT_FILE       = "server.crt"
	KEY_FILE        = "server.key"
	DATE_FORMAT     = "2006-01-02"
	USERNAME        = "user"
	PASSWORD        = "testpassword"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

type flags struct {
	orderId         uint64
	customerId      uint64
	keepDateString  string
	numberOfEntries int
	storedOnly      bool
	page            int
}

func initArgs() (cmd string, args []string, f flags, err error) {
	fs := flag.NewFlagSet("", flag.ExitOnError)
	fs.Usage = printHelp
	fs.Uint64Var(&f.orderId, "order-id", 0, "specify order id")
	fs.Uint64Var(&f.customerId, "customer-id", 0, "specify customer id")
	fs.StringVar(&f.keepDateString, "keep-date", "", "specify keep date")
	fs.IntVar(&f.numberOfEntries, "n", 0, "specify number of entries")
	fs.BoolVar(&f.storedOnly, "stored-only", false, "display only stored orders")
	fs.IntVar(&f.page, "page", 0, "specify page")

	if len(os.Args) < 2 {
		return "", nil, flags{}, errors.New("subcommand is required, see `help` subcommand for details")
	}
	cmd = os.Args[1]
	args = os.Args[2:]
	err = fs.Parse(args)
	if err != nil {
		return "", nil, flags{}, err
	}
	return cmd, args, f, nil
}

func run() error {
	cmd, args, f, err := initArgs()
	if err != nil {
		return err
	}

	stor, err := file.NewOrderFileStorage(ORDERS_FILEPATH)
	if err != nil {
		return err
	}
	serv := order.NewOrderService(&stor)
	defer stor.Close()

	switch cmd {
	case "help":
		printHelp()
		return nil
	case "manage-pickup-points":
		return managePickUpPoints()
	case "run-pickup-points-api":
		return runPickUpPointRestApi()
	case "accept-order":
		return acceptOrder(f, serv)
	case "return-order":
		return returnOrder(f, serv)
	case "give-orders":
		return giveOrders(args, serv)
	case "list-orders":
		return listOrders(f, serv)
	case "accept-return":
		return acceptReturn(f, serv)
	case "list-returns":
		return listReturns(serv, f)
	default:
		return errors.New("not such subcommand, see `help` subcommand for details")
	}
}

func printHelp() {
	fmt.Fprintln(os.Stderr, `Available commands:

	help
		Show this help message

	manage-pickup-points
		Starts interactive mode for managing pick-up points

	run-pickup-points-api
		Starts a HTTPS API server for managing pick-up points

	accept-order --order-id <order-id> --customer-id <customer-id> --keep-date <keep-date>
		Accepts order from a courier
		--order-id		specify an order id
		--customer-id	specify a customer id
		--keep-date		specify a keep date in YYYY-MM-DD format

	return-order --order-id <order-id>
		Returns order to a courier
		--order-id		specify an order id

	give-orders <order-id> ...
		Gives orders to a customer

	list-orders --customer-id <customer-id> [-n <number-of-entries>] [--stored-only]
		Lists registered orders
		--customer-id	specify a customer id
		-n				specify max number of entries to display
		--stored-only	show only orders which are currently stored

	accept-return --customer-id <customer-id> --order-id <order-id>
		Accepts order return from a customer
		--customer-id	specify a customer id
		--order-id		specify an order id

	list-returns -n <number-of-entries> [--page <page-num>]
		Lists all stored returned orders
		-n				specify number of entries to display on a page
		--page			specify page number, starting with 0`)
}

func managePickUpPoints() error {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	stor, err := file.NewPickUpPointFileStorage(POINTS_FILEPATH)
	if err != nil {
		return err
	}
	defer stor.Close()

	eg, ctx := errgroup.WithContext(ctx)

	logCtx, stopLog := context.WithCancel(context.Background())
	defer stopLog()

	thrLog := logger.NewLogger()
	go thrLog.Run(logCtx)

	serv := pickuppoint.NewPickUpPointService(&stor, thrLog)
	eg.Go(func() error {
		return serv.Run(ctx)
	})

	pickUpPointCli := cli.NewPickUpPointCli(serv, thrLog)
	eg.Go(func() error {
		return pickUpPointCli.Run(ctx)
	})

	return eg.Wait()
}

func runPickUpPointRestApi() error {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	pointDb, err := db.NewDb(ctx)
	if err != nil {
		return err
	}
	defer pointDb.Close()

	stor := postgres.NewPickUpPointStorage(pointDb)

	eg, ctx := errgroup.WithContext(ctx)

	logCtx, stopLog := context.WithCancel(context.Background())
	defer stopLog()

	thrLog := logger.NewLogger()
	go thrLog.Run(logCtx)

	serv := httpserv.NewHttpServer(stor, thrLog)
	eg.Go(func() error {
		return serv.Serve(ctx, HTTPS_ADDR, REDIRECT_ADDR, CERT_FILE, KEY_FILE, USERNAME, PASSWORD)
	})

	return eg.Wait()
}

func acceptOrder(f flags, serv order.Service) error {
	if f.keepDateString == "" {
		return errors.New("keep date is required")
	}
	keepDate, err := time.ParseInLocation(DATE_FORMAT, f.keepDateString, time.Local)
	if err != nil {
		return err
	}
	keepDate = keepDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	return serv.AddOrder(f.orderId, f.customerId, keepDate)
}

func returnOrder(f flags, serv order.Service) error {
	return serv.RemoveOrder(f.orderId)
}

func giveOrders(args []string, serv order.Service) error {
	orderIds := make([]uint64, len(args))
	for i, arg := range args {
		var err error
		orderIds[i], err = strconv.ParseUint(arg, 10, 64)
		if err != nil {
			return err
		}
	}
	return serv.GiveOrders(orderIds)
}

func listOrders(f flags, serv order.Service) error {
	orders, err := serv.GetOrders(f.customerId, f.numberOfEntries, f.storedOnly)
	if err != nil {
		return err
	}
	if len(orders) == 0 {
		fmt.Println("No orders are to display")
		return nil
	}
	printOrders(orders)
	return nil
}

func acceptReturn(f flags, serv order.Service) error {
	return serv.AcceptReturn(f.orderId, f.customerId)
}

func listReturns(serv order.Service, f flags) error {
	orders, err := serv.GetReturns(f.numberOfEntries, f.page)
	if err != nil {
		return err
	}
	if len(orders) == 0 {
		fmt.Println("No returns are available")
		return nil
	}
	printOrders(orders)
	return nil
}

func printOrders(orders []model.Order) {
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintf(
		w,
		"%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
		"Order id",
		"Customer id",
		"Add date",
		"Keep date",
		"Is given",
		"Give date",
		"Is returned",
		"Return date")
	for _, order := range orders {
		fmt.Fprintln(w, order)
	}
	w.Flush()
}
