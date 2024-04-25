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

type PickUpPointGrpcConsoleCommands struct {
	svc  core.PickUpPointCoreService
	help Command
}

func NewPickUpPointGrpcConsoleCommands(svc core.PickUpPointCoreService, help Command) *PickUpPointGrpcConsoleCommands {
	return &PickUpPointGrpcConsoleCommands{svc: svc, help: help}
}

func (c *PickUpPointGrpcConsoleCommands) RunPickUpPointGrpcApi(args []string) error {
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

	handler := grpchandler.NewService(c.svc, log)

	serv := grpc.NewServer()
	pb.RegisterPickUpPointServiceServer(serv, handler)

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return serv.Serve(lis)
	})

	return eg.Wait()
}
