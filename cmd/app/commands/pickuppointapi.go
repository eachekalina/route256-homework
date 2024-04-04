package commands

import (
	"context"
	"github.com/gorilla/mux"
	"golang.org/x/sync/errgroup"
	"homework/cmd/app/httpserv"
	"homework/internal/app/core"
	"homework/internal/app/logger"
	"homework/internal/app/middleware"
	"net/http"
	"os/signal"
	"syscall"
)

type PickUpPointApiConsoleCommands struct {
	svc  core.PickUpPointCoreService
	help Command
}

func NewPickUpPointApiConsoleCommands(svc core.PickUpPointCoreService, help Command) *PickUpPointApiConsoleCommands {
	return &PickUpPointApiConsoleCommands{svc: svc, help: help}
}

func (c *PickUpPointApiConsoleCommands) RunPickUpPointApi(args []string) error {
	var params httpserv.HttpServerParams
	var username, password string

	fs := createFlagSet(c.help)
	fs.StringVar(&params.HttpsAddr, "https-address", ":9443", "specify https listen address")
	fs.StringVar(&params.RedirectAddr, "redirect-address", ":9000", "specify redirect listen address")
	fs.StringVar(&params.CertFile, "tls-cert", "server.crt", "specify tls certificate file")
	fs.StringVar(&params.KeyFile, "tls-key", "server.key", "specify tls certificate key file")
	fs.StringVar(&username, "username", "user", "specify access control username")
	fs.StringVar(&password, "password", "testpassword", "specify access control password")
	err := fs.Parse(args)
	if err != nil {
		return err
	}

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	eg, ctx := errgroup.WithContext(ctx)

	logCtx, stopLog := context.WithCancel(context.Background())
	defer stopLog()
	log := logger.NewLogger()
	go log.Run(logCtx)

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
		middleware.LogMiddleware(log),
		middleware.AuthMiddleware(username, password),
	}

	serv := httpserv.NewHttpServer(params)
	eg.Go(func() error {
		return serv.Serve(ctx)
	})

	return eg.Wait()
}
