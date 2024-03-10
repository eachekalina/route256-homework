package main

import (
	"Homework-1/internal/model"
	"Homework-1/internal/service"
	"Homework-1/internal/storage"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"text/tabwriter"
	"time"
)

const FILEPATH = "orders.json"

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

const dateFormat = "2006-01-02"

func run() error {
	cmd, args, f, err := initArgs()
	if err != nil {
		return err
	}

	stor, err := storage.NewFileStorage(FILEPATH)
	if err != nil {
		return err
	}
	serv := service.NewService(&stor)
	defer stor.Close()

	switch cmd {
	case "help":
		printHelp()
		return nil
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

func acceptOrder(f flags, serv service.Service) error {
	if f.keepDateString == "" {
		return errors.New("keep date is required")
	}
	keepDate, err := time.Parse(dateFormat, f.keepDateString)
	if err != nil {
		return err
	}
	return serv.AddOrder(f.orderId, f.customerId, keepDate)
}

func returnOrder(f flags, serv service.Service) error {
	return serv.RemoveOrder(f.orderId)
}

func giveOrders(args []string, serv service.Service) error {
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

func listOrders(f flags, serv service.Service) error {
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

func acceptReturn(f flags, serv service.Service) error {
	return serv.AcceptReturn(f.orderId, f.customerId)
}

func listReturns(serv service.Service, f flags) error {
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
