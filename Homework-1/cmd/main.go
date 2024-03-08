package main

import (
	"Homework-1/internal/model"
	"Homework-1/internal/service"
	"Homework-1/internal/storage"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"
)

const FILEPATH = "orders.json"

func main() {
	err := run()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	fs := flag.NewFlagSet("main", flag.ExitOnError)
	orderId := fs.Uint64("order-id", 0, "specify order id")
	customerId := fs.Uint64("customer-id", 0, "specify customer id")
	keepDateString := fs.String("keep-date", "", "specify keep date")
	numberOfEntries := fs.Int("n", 0, "specify number of entries")
	storedOnly := fs.Bool("stored-only", false, "display only stored orders")
	page := fs.Int("page", 0, "specify page")

	if len(os.Args) < 2 {
		return errors.New("subcommand is required")
	}
	cmd := os.Args[1]
	args := os.Args[2:]
	err := fs.Parse(args)
	if err != nil {
		return err
	}

	stor, err := storage.NewFileStorage(FILEPATH)
	if err != nil {
		return err
	}
	serv := service.New(&stor)

	switch cmd {
	case "accept-order":
		if *orderId == 0 {
			return errors.New("valid order id is required")
		}
		if *customerId == 0 {
			return errors.New("valid customer id is required")
		}
		if *keepDateString == "" {
			return errors.New("keep date is required")
		}
		keepDate, err := time.Parse("2006-01-02", *keepDateString)
		if err != nil {
			return err
		}
		err = serv.AddOrder(*orderId, *customerId, keepDate)
		if err != nil {
			return err
		}
	case "return-order":
		if *orderId == 0 {
			return errors.New("valid order id is required")
		}
		err := serv.RemoveOrder(*orderId)
		if err != nil {
			return err
		}
	case "give-orders":
		orderIds := make([]uint64, len(args))
		for i, arg := range args {
			orderIds[i], err = strconv.ParseUint(arg, 10, 64)
			if err != nil {
				return err
			}
		}
		err := serv.GiveOrders(orderIds)
		if err != nil {
			return err
		}
	case "list-orders":
		if *customerId == 0 {
			return errors.New("valid customer id is required")
		}
		if *numberOfEntries == 0 {
			return errors.New("valid number of entries is required")
		}
		orders, err := serv.GetOrders(*customerId, *numberOfEntries, *storedOnly)
		if err != nil {
			return err
		}
		if len(orders) == 0 {
			fmt.Println("No orders are to display")
			break
		}
		printOrders(orders)
	case "accept-return":
		if *orderId == 0 {
			return errors.New("valid order id is required")
		}
		if *customerId == 0 {
			return errors.New("valid customer id is required")
		}
		err := serv.AcceptReturn(*orderId, *customerId)
		if err != nil {
			return err
		}
	case "list-returns":
		orders, err := serv.GetReturns(*numberOfEntries, *page)
		if err != nil {
			return err
		}
		if len(orders) == 0 {
			fmt.Println("No returns are available")
			break
		}
		printOrders(orders)
	}
	err = stor.Save()
	if err != nil {
		return err
	}
	return nil
}

func printOrders(orders []model.Order) {
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	_, _ = fmt.Fprintf(
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
		displayedGiveDate := "-"
		if order.IsGiven {
			displayedGiveDate = order.GiveDate.Format("2006-01-02")
		}
		displayedReturnDate := "-"
		if order.IsReturned {
			displayedReturnDate = order.ReturnDate.Format("2006-01-02")
		}
		_, _ = fmt.Fprintf(
			w,
			"%d\t%d\t%s\t%s\t%t\t%s\t%t\t%s\n",
			order.Id,
			order.CustomerId,
			order.AddDate.Format("2006-01-02"),
			order.KeepDate.Format("2006-01-02"),
			order.IsGiven,
			displayedGiveDate,
			order.IsReturned,
			displayedReturnDate)
	}
	_ = w.Flush()
}
