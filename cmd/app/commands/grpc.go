package commands

import (
	"context"
	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"homework/cmd/app/grpchandler"
	"homework/internal/app/core"
	"homework/internal/app/logger"
	"homework/internal/app/metrics"
	"homework/internal/app/pb"
	"net"
	"net/http"
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
	var tracingAddr string

	fs := createFlagSet(c.help)
	fs.StringVar(&listenAddr, "listen-address", ":9090", "specify listen address")
	fs.StringVar(&tracingAddr, "tracing-address", "localhost:4318", "specify tracer host address")
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

	reg := prometheus.NewRegistry()

	grpcMetrics := grpcprometheus.NewServerMetrics()
	reg.MustRegister(grpcMetrics)

	metric := metrics.NewMetrics(reg)

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	eg, ctx := errgroup.WithContext(ctx)

	shutdown, err := metrics.NewProvider(ctx, tracingAddr)
	if err != nil {
		return err
	}
	defer shutdown(ctx)

	pointHandler := grpchandler.NewPickUpPointService(c.pointSvc, log)
	orderHandler := grpchandler.NewOrderService(c.orderSvc, log, metric)

	serv := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(grpcMetrics.UnaryServerInterceptor()),
		grpc.ChainStreamInterceptor(grpcMetrics.StreamServerInterceptor()),
	)
	pb.RegisterPickUpPointServiceServer(serv, pointHandler)
	pb.RegisterOrderServiceServer(serv, orderHandler)

	eg.Go(func() error {
		return serv.Serve(lis)
	})

	go http.ListenAndServe(":9091", promhttp.HandlerFor(reg, promhttp.HandlerOpts{EnableOpenMetrics: true}))

	<-ctx.Done()
	serv.GracefulStop()

	return eg.Wait()
}
