package httpserv

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"golang.org/x/sync/errgroup"
	"homework/internal/model"
	storerr "homework/internal/storage"
	"io"
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
	stor storage
}

func NewHttpServer(stor storage) *HttpServer {
	return &HttpServer{stor: stor}
}

func (s *HttpServer) Serve(ctx context.Context) {
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

	httpsServer := http.Server{}
	redirectServer := http.Server{}

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(httpsServer.ListenAndServe)
	eg.Go(redirectServer.ListenAndServe)
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
