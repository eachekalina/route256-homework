package main

import (
	"context"
	"fmt"
	"homework/cmd/app/commands"
	"homework/internal/app/core"
	"homework/internal/app/db"
	"homework/internal/app/order"
	"homework/internal/app/packaging"
	"homework/internal/app/pickuppoint"
	"log"
	"os"
)

const (
	ORDERS_FILEPATH = "orders.json"
	POINTS_FILEPATH = "points.json"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	helpCommand := func(args []string) error {
		help()
		return nil
	}

	pointFileRepo, err := pickuppoint.NewFileRepository(POINTS_FILEPATH)
	if err != nil {
		return err
	}

	cliCommands := commands.NewPickUpPointCliConsoleCommands(
		core.NewPickUpPointCoreService(pickuppoint.NewService(pointFileRepo)),
		helpCommand,
	)

	pointDb, err := db.NewDb(context.Background())
	if err != nil {
		return err
	}

	apiCommands := commands.NewPickUpPointApiConsoleCommands(
		core.NewPickUpPointCoreService(pickuppoint.NewService(pickuppoint.NewPostgresRepository(pointDb))),
		helpCommand,
	)

	orderFileRepo, err := order.NewOrderFileRepository(ORDERS_FILEPATH)
	if err != nil {
		return err
	}

	packagingTypes := map[packaging.Type]packaging.Packaging{
		packaging.BagType:  packaging.Bag{},
		packaging.BoxType:  packaging.Box{},
		packaging.FilmType: packaging.Film{},
	}

	orderCommands := commands.NewOrderConsoleCommands(
		core.NewOrderCoreService(order.NewService(orderFileRepo), packagingTypes),
		helpCommand,
	)

	cmdMap := map[string]commands.Command{
		"help":                  helpCommand,
		"manage-pickup-points":  cliCommands.ManagePickUpPointsCommand,
		"run-pickup-points-api": apiCommands.RunPickUpPointApi,
		"accept-order":          orderCommands.AcceptOrderCommand,
		"return-order":          orderCommands.ReturnOrderCommand,
		"give-orders":           orderCommands.GiveOrdersCommand,
		"list-orders":           orderCommands.ListOrdersCommand,
		"accept-return":         orderCommands.AcceptReturnCommand,
		"list-returns":          orderCommands.ListReturnsCommand,
	}
	return commands.Run(cmdMap)
}

func help() {
	fmt.Fprintln(os.Stderr, `Available commands:

	help
		Show this help message

	manage-pickup-points
		Starts interactive mode for managing pick-up points

	run-pickup-points-api
		Starts a HTTPS API server for managing pick-up points
		--https-address		specify HTTPS listen address, default: :9443
		--redirect-address	specify redirect listen address, default: :9000
		--tls-cert			specify TLS certificate file, default: server.crt
		--tls-key			specify TLS certificate key file, default: server.key
		--username			specify access control username, default: user
		--password			specify access control password, default: testpassword

	accept-order --order-id <order-id> --customer-id <customer-id> --keep-date <keep-date> --price <price> --weight <weight>
		Accepts order from a courier
		--order-id		specify an order id
		--customer-id	specify a customer id
		--keep-date		specify a keep date in YYYY-MM-DD format
		--price			specify price in rubles
		--weight		specify weight in kg
		--packaging		specify packaging type

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
