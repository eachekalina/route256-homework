package httpserv

import (
	"context"
	"github.com/gorilla/mux"
	"golang.org/x/sync/errgroup"
	"net"
	"net/http"
)

type Handler func(req *http.Request) (int, []byte)

type PathHandler struct {
	Methods map[string]Handler
}

type HttpServerParams struct {
	Handlers     map[string]PathHandler
	Middlewares  []mux.MiddlewareFunc
	HttpsAddr    string
	RedirectAddr string
	CertFile     string
	KeyFile      string
	Username     string
	Password     string
}

type HttpServer struct {
	params HttpServerParams
}

func NewHttpServer(params HttpServerParams) *HttpServer {
	return &HttpServer{params: params}
}

func (s *HttpServer) Serve(ctx context.Context) error {
	router := mux.NewRouter()
	for path, pathHandler := range s.params.Handlers {
		router.HandleFunc(path, s.makeHandlerFunc(pathHandler))
	}
	for _, middleware := range s.params.Middlewares {
		router.Use(middleware)
	}

	_, httpsPort, err := net.SplitHostPort(s.params.HttpsAddr)
	if err != nil {
		return err
	}

	httpsServer := http.Server{
		Addr:    s.params.HttpsAddr,
		Handler: router,
	}
	redirectServer := http.Server{
		Addr: s.params.RedirectAddr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			host, _, err := net.SplitHostPort(req.Host)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			u := *req.URL
			u.Host = net.JoinHostPort(host, httpsPort)
			u.Scheme = "https"
			http.Redirect(w, req, u.String(), http.StatusMovedPermanently)
		}),
	}

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return httpsServer.ListenAndServeTLS(s.params.CertFile, s.params.KeyFile)
	})
	eg.Go(redirectServer.ListenAndServe)

	eg.Go(func() error {
		<-ctx.Done()
		redirectServer.Shutdown(context.Background())
		httpsServer.Shutdown(context.Background())
		return nil
	})

	return eg.Wait()
}

func (s *HttpServer) makeHandlerFunc(pathHandler PathHandler) http.HandlerFunc {
	var methods string
	for method := range pathHandler.Methods {
		if methods == "" {
			methods = method
			continue
		}
		methods = methods + ", " + method
	}
	return func(w http.ResponseWriter, req *http.Request) {
		var code int
		var body []byte
		handler, ok := pathHandler.Methods[req.Method]
		if ok {
			code, body = handler(req)
		} else {
			w.Header().Set("Allow", methods)
			code = http.StatusMethodNotAllowed
		}
		w.WriteHeader(code)
		if body != nil {
			w.Write(body)
		}
	}
}
