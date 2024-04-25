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
	filePerm        = 0777
	topic           = "requests"
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

	pointFileRepo, closePointFileRepo, err := initPickUpPointFileRepository(POINTS_FILEPATH, filePerm)
	if err != nil {
		return err
	}
	defer closePointFileRepo()

	cliCommands := commands.NewPickUpPointCliConsoleCommands(
		core.NewPickUpPointCoreService(pickuppoint.NewService(pointFileRepo, db.Dummy{})),
		helpCommand,
	)

	tm, err := db.NewTransactionManager(context.Background())
	if err != nil {
		return err
	}

	apiCoreService := core.NewPickUpPointCoreService(pickuppoint.NewService(pickuppoint.NewPostgresRepository(db.NewDatabase(tm)), tm))

	apiCommands := commands.NewPickUpPointApiConsoleCommands(apiCoreService, helpCommand, topic)

	grpcCommands := commands.NewPickUpPointGrpcConsoleCommands(apiCoreService, helpCommand)

	orderFileRepo, closeOrderFileRepo, err := initOrderFileRepository(ORDERS_FILEPATH, filePerm)
	if err != nil {
		return err
	}
	defer closeOrderFileRepo()

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
		"help":                   helpCommand,
		"manage-pickup-points":   cliCommands.ManagePickUpPointsCommand,
		"run-pickup-points-api":  apiCommands.RunPickUpPointApi,
		"run-pickup-points-grpc": grpcCommands.RunPickUpPointGrpcApi,
		"accept-order":           orderCommands.AcceptOrderCommand,
		"return-order":           orderCommands.ReturnOrderCommand,
		"give-orders":            orderCommands.GiveOrdersCommand,
		"list-orders":            orderCommands.ListOrdersCommand,
		"accept-return":          orderCommands.AcceptReturnCommand,
		"list-returns":           orderCommands.ListReturnsCommand,
	}
	return commands.Run(cmdMap)
}

func initPickUpPointFileRepository(pointsFilePath string, filePerm os.FileMode) (*pickuppoint.FileRepository, func(), error) {
	file, err := os.OpenFile(pointsFilePath, os.O_CREATE|os.O_RDONLY, filePerm)
	if err != nil {
		return nil, nil, err
	}
	repo, err := pickuppoint.NewFileRepository(file)
	if err != nil {
		closeErr := file.Close()
		if closeErr != nil {
			log.Println(closeErr)
		}
		return nil, nil, err
	}
	err = file.Close()
	if err != nil {
		return nil, nil, err
	}
	f := func() {
		file, err := os.OpenFile(POINTS_FILEPATH, os.O_CREATE|os.O_WRONLY, filePerm)
		if err != nil {
			log.Println(err)
		}
		err = repo.Close(file)
		if err != nil {
			log.Println(err)
		}
	}
	return repo, f, nil
}

func initOrderFileRepository(orderFilePath string, filePerm os.FileMode) (*order.FileRepository, func(), error) {
	file, err := os.OpenFile(orderFilePath, os.O_CREATE|os.O_RDONLY, filePerm)
	if err != nil {
		return nil, nil, err
	}
	repo, err := order.NewFileRepository(file)
	if err != nil {
		closeErr := file.Close()
		if closeErr != nil {
			log.Println(closeErr)
		}
		return nil, nil, err
	}
	err = file.Close()
	if err != nil {
		return nil, nil, err
	}
	f := func() {
		file, err := os.OpenFile(ORDERS_FILEPATH, os.O_CREATE|os.O_WRONLY, filePerm)
		if err != nil {
			log.Println(err)
		}
		err = repo.Close(file)
		if err != nil {
			log.Println(err)
		}
	}
	return repo, f, nil
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
		--brokers			specify broker addresses, separated by comma, default: 127.0.0.1:9091,127.0.0.1:9092,127.0.0.1:9093

	run-pickup-points-grpc
		Starts a gRPC API server for managing pick-up points
		--listen-address	specify listen address, default: :9090

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
