package commands

import (
	"context"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"homework/cmd/app/grpchandler"
	"homework/internal/app/core"
	"homework/internal/app/logger"
	"homework/internal/app/pb"
	"net"
	"os/signal"
	"syscall"
)

type GrpcConsoleCommands struct {
	pointSvc core.PickUpPointCoreService
	orderSvc *core.OrderCoreService
	help     Command
}

func NewGrpcConsoleCommands(pointSvc core.PickUpPointCoreService, orderSvc *core.OrderCoreService, help Command) *GrpcConsoleCommands {
	return &GrpcConsoleCommands{pointSvc: pointSvc, orderSvc: orderSvc, help: help}
}

func (c *GrpcConsoleCommands) RunGrpcApi(args []string) error {
	var listenAddr string

	fs := createFlagSet(c.help)
	fs.StringVar(&listenAddr, "listen-address", ":9090", "specify listen address")
	err := fs.Parse(args)
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}

	logCtx, stopLog := context.WithCancel(context.Background())
	defer stopLog()
	log := logger.NewLogger()
	go log.Run(logCtx)

	pointHandler := grpchandler.NewPickUpPointService(c.pointSvc, log)
	orderHandler := grpchandler.NewOrderService(c.orderSvc, log)

	serv := grpc.NewServer()
	pb.RegisterPickUpPointServiceServer(serv, pointHandler)
	pb.RegisterOrderServiceServer(serv, orderHandler)

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return serv.Serve(lis)
	})

	<-ctx.Done()
	serv.GracefulStop()

	return eg.Wait()
}
