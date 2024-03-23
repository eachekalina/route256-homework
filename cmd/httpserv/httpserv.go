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
		switch req.Method {
		case http.MethodGet:
			s.list(w, req)
		case http.MethodPost:
			s.create(w, req)
		default:
			w.Header().Set("Allow", "GET, POST")
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	router.HandleFunc("/pickup-point/{id:[0-9]+}", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			s.get(w, req)
		case http.MethodPut:
			s.update(w, req)
		case http.MethodDelete:
			s.delete(w, req)
		default:
			w.Header().Set("Allow", "GET, PUT, DELETE")
			w.WriteHeader(http.StatusMethodNotAllowed)
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

func (s *HttpServer) create(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var point model.PickUpPoint
	err = json.Unmarshal(body, &point)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.stor.Create(req.Context(), point)
	if err != nil {
		if errors.Is(err, storerr.ErrIdAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pointJson, err := json.Marshal(point)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(pointJson)
}

func (s *HttpServer) list(w http.ResponseWriter, req *http.Request) {
	list, err := s.stor.List(req.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	listJson, err := json.Marshal(list)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(listJson)
}

func (s *HttpServer) get(w http.ResponseWriter, req *http.Request) {
	idStr, ok := mux.Vars(req)["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	point, err := s.stor.Get(req.Context(), id)
	if err != nil {
		if errors.Is(err, storerr.ErrNoItemFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pointJson, err := json.Marshal(point)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(pointJson)
}

func (s *HttpServer) update(w http.ResponseWriter, req *http.Request) {
	idStr, ok := mux.Vars(req)["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var point model.PickUpPoint
	err = json.Unmarshal(body, &point)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if point.Id != id {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.stor.Update(req.Context(), point)
	if err != nil {
		if errors.Is(err, storerr.ErrNoItemFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pointJson, err := json.Marshal(point)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(pointJson)
}

func (s *HttpServer) delete(w http.ResponseWriter, req *http.Request) {
	idStr, ok := mux.Vars(req)["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.stor.Delete(req.Context(), id)
	if err != nil {
		if errors.Is(err, storerr.ErrNoItemFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
