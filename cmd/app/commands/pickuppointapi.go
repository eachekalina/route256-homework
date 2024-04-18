package commands

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"golang.org/x/sync/errgroup"
	"homework/cmd/app/httpserv"
	"homework/internal/app/core"
	"homework/internal/app/kafka"
	"homework/internal/app/logger"
	"homework/internal/app/middleware"
	"homework/internal/app/reqlog"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type PickUpPointApiConsoleCommands struct {
	svc   core.PickUpPointCoreService
	help  Command
	topic string
}

func NewPickUpPointApiConsoleCommands(svc core.PickUpPointCoreService, help Command, topic string) *PickUpPointApiConsoleCommands {
	return &PickUpPointApiConsoleCommands{svc: svc, help: help, topic: topic}
}

func (c *PickUpPointApiConsoleCommands) RunPickUpPointApi(args []string) error {
	var params httpserv.HttpServerParams
	var username, password string
	var brokersStr string

	fs := createFlagSet(c.help)
	fs.StringVar(&params.HttpsAddr, "https-address", ":9443", "specify https listen address")
	fs.StringVar(&params.RedirectAddr, "redirect-address", ":9000", "specify redirect listen address")
	fs.StringVar(&params.CertFile, "tls-cert", "server.crt", "specify tls certificate file")
	fs.StringVar(&params.KeyFile, "tls-key", "server.key", "specify tls certificate key file")
	fs.StringVar(&username, "username", "user", "specify access control username")
	fs.StringVar(&password, "password", "testpassword", "specify access control password")
	fs.StringVar(&brokersStr, "brokers", "127.0.0.1:9091,127.0.0.1:9092,127.0.0.1:9093", "specify broker addresses, separated by comma")
	err := fs.Parse(args)
	if err != nil {
		return err
	}

	brokers := strings.Split(brokersStr, ",")

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	eg, ctx := errgroup.WithContext(ctx)

	logCtx, stopLog := context.WithCancel(context.Background())
	defer stopLog()
	log := logger.NewLogger()
	go log.Run(logCtx)

	producer, err := kafka.NewProducer(brokers, log, c.topic)
	if err != nil {
		return err
	}
	defer producer.Close()
	consumer, err := kafka.NewConsumer(brokers, c.topic, reqlog.LogHandler(log))
	if err != nil {
		return err
	}
	eg.Go(func() error {
		return consumer.Run(ctx)
	})
	t := time.NewTicker(5 * time.Second)
	select {
	case <-consumer.Ready():
		log.Log("Consumer ready")
	case <-t.C:
		t.Stop()
		return errors.New("consumer failed: timeout")
	}
	t.Stop()
	reqLog := reqlog.NewLogger(producer, consumer)

	handlers := httpserv.NewPickUpPointHandlers(c.svc, log)

	params.Handlers = map[string]httpserv.PathHandler{
		"/pickup-point": {Methods: map[string]httpserv.Handler{
			http.MethodGet:  handlers.ListHandler,
			http.MethodPost: handlers.CreateHandler,
		}},
		"/pickup-point/{id:[0-9]+}": {Methods: map[string]httpserv.Handler{
			http.MethodGet:    handlers.GetHandler,
			http.MethodPut:    handlers.UpdateHandler,
			http.MethodDelete: handlers.DeleteHandler,
		}},
	}

	params.Middlewares = []mux.MiddlewareFunc{
		middleware.LogMiddleware(reqLog),
		middleware.AuthMiddleware(username, password),
	}

	serv := httpserv.NewHttpServer(params)
	eg.Go(func() error {
		return serv.Serve(ctx)
	})

	return eg.Wait()
}
