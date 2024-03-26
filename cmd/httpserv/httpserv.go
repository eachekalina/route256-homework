package httpserv

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"golang.org/x/sync/errgroup"
	"homework/internal/logger"
	"homework/internal/model"
	storerr "homework/internal/storage"
	"io"
	"net"
	"net/http"
	"strconv"
)

type storage interface {
	Create(ctx context.Context, point model.PickUpPoint) error
	List(ctx context.Context) ([]model.PickUpPoint, error)
	Get(ctx context.Context, id uint64) (model.PickUpPoint, error)
	Update(ctx context.Context, point model.PickUpPoint) error
	Delete(ctx context.Context, id uint64) error
}

type HttpServer struct {
	stor     storage
	log      *logger.Logger
	username string
	password string
}

func NewHttpServer(stor storage, log *logger.Logger) *HttpServer {
	return &HttpServer{stor: stor, log: log}
}

func (s *HttpServer) Serve(ctx context.Context, httpsAddr string, redirectAddr string, certFile string, keyFile string, username string, password string) error {
	router := mux.NewRouter()
	router.HandleFunc("/pickup-point", func(w http.ResponseWriter, req *http.Request) {
		var code int
		var body []byte
		switch req.Method {
		case http.MethodGet:
			code, body = s.list(req)
		case http.MethodPost:
			code, body = s.create(req)
		default:
			w.Header().Set("Allow", "GET, POST")
			code = http.StatusMethodNotAllowed
		}
		w.WriteHeader(code)
		if body != nil {
			w.Write(body)
		}
	})
	router.HandleFunc("/pickup-point/{id:[0-9]+}", func(w http.ResponseWriter, req *http.Request) {
		var code int
		var body []byte
		switch req.Method {
		case http.MethodGet:
			code, body = s.get(req)
		case http.MethodPut:
			code, body = s.update(req)
		case http.MethodDelete:
			code, body = s.delete(req)
		default:
			w.Header().Set("Allow", "GET, PUT, DELETE")
			code = http.StatusMethodNotAllowed
		}
		w.WriteHeader(code)
		if body != nil {
			w.Write(body)
		}
	})
	router.Use(s.logMiddleware)
	router.Use(s.authMiddleware)

	s.username = username
	s.password = password

	_, httpsPort, err := net.SplitHostPort(httpsAddr)
	if err != nil {
		return err
	}

	httpsServer := http.Server{
		Addr:    httpsAddr,
		Handler: router,
	}
	redirectServer := http.Server{
		Addr: redirectAddr,
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
		return httpsServer.ListenAndServeTLS(certFile, keyFile)
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

func (s *HttpServer) create(req *http.Request) (int, []byte) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return http.StatusInternalServerError, nil
	}

	var point model.PickUpPoint
	err = json.Unmarshal(body, &point)
	if err != nil {
		return http.StatusBadRequest, nil
	}

	err = s.stor.Create(req.Context(), point)
	if err != nil {
		if errors.Is(err, storerr.ErrIdAlreadyExists) {
			return http.StatusConflict, nil
		}
		return http.StatusInternalServerError, nil
	}

	pointJson, err := json.Marshal(point)
	if err != nil {
		return http.StatusInternalServerError, nil
	}

	return http.StatusCreated, pointJson
}

func (s *HttpServer) list(req *http.Request) (int, []byte) {
	list, err := s.stor.List(req.Context())
	if err != nil {
		return http.StatusInternalServerError, nil
	}

	listJson, err := json.Marshal(list)
	if err != nil {
		return http.StatusInternalServerError, nil
	}

	return http.StatusOK, listJson
}

func (s *HttpServer) get(req *http.Request) (int, []byte) {
	idStr, ok := mux.Vars(req)["id"]
	if !ok {
		return http.StatusBadRequest, nil
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return http.StatusBadRequest, nil
	}

	point, err := s.stor.Get(req.Context(), id)
	if err != nil {
		if errors.Is(err, storerr.ErrNoItemFound) {
			return http.StatusNotFound, nil
		}
		return http.StatusInternalServerError, nil
	}

	pointJson, err := json.Marshal(point)
	if err != nil {
		return http.StatusInternalServerError, nil
	}

	return http.StatusOK, pointJson
}

func (s *HttpServer) update(req *http.Request) (int, []byte) {
	idStr, ok := mux.Vars(req)["id"]
	if !ok {
		return http.StatusBadRequest, nil
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return http.StatusBadRequest, nil
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return http.StatusInternalServerError, nil
	}

	var point model.PickUpPoint
	err = json.Unmarshal(body, &point)
	if err != nil {
		return http.StatusBadRequest, nil
	}

	if point.Id != id {
		return http.StatusBadRequest, nil
	}

	err = s.stor.Update(req.Context(), point)
	if err != nil {
		if errors.Is(err, storerr.ErrNoItemFound) {
			return http.StatusNotFound, nil
		}
		return http.StatusInternalServerError, nil
	}

	pointJson, err := json.Marshal(point)
	if err != nil {
		return http.StatusInternalServerError, nil
	}

	return http.StatusOK, pointJson
}

func (s *HttpServer) delete(req *http.Request) (int, []byte) {
	idStr, ok := mux.Vars(req)["id"]
	if !ok {
		return http.StatusBadRequest, nil
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return http.StatusBadRequest, nil
	}

	err = s.stor.Delete(req.Context(), id)
	if err != nil {
		if errors.Is(err, storerr.ErrNoItemFound) {
			return http.StatusNotFound, nil
		}
		return http.StatusInternalServerError, nil
	}

	return http.StatusNoContent, nil
}

func (s *HttpServer) logMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var buf bytes.Buffer
		req.Body = io.NopCloser(io.TeeReader(req.Body, &buf))
		h.ServeHTTP(w, req)
		s.log.Log("Request from %s, method %s, path %s, headers %s, body %s", req.RemoteAddr, req.Method, req.RequestURI, req.Header, buf.String())
	})
}

func (s *HttpServer) authMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		username, password, ok := req.BasicAuth()
		if !ok || username != s.username || password != s.password {
			w.Header().Set("WWW-Authenticate", "Basic realm=\"pickup-point\"")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, req)
	})
}
