package httpserv

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"homework/internal/app/core"
	"homework/internal/app/logger"
	"homework/internal/app/pickuppoint"
	"io"
	"net/http"
	"strconv"
)

type PickUpPointHandlers struct {
	svc *core.PickUpPointCoreService
	log *logger.Logger
}

func NewPickUpPointHandlers(svc *core.PickUpPointCoreService, log *logger.Logger) *PickUpPointHandlers {
	return &PickUpPointHandlers{svc: svc, log: log}
}

func (h *PickUpPointHandlers) CreateHandler(httpReq *http.Request) (int, []byte) {
	body, err := io.ReadAll(httpReq.Body)
	if err != nil {
		h.log.Log("%v", err)
		return http.StatusInternalServerError, nil
	}

	var req core.CreatePointRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		return http.StatusBadRequest, nil
	}

	err = h.svc.CreatePoint(httpReq.Context(), req)
	if err != nil {
		if errors.Is(err, pickuppoint.ErrIdAlreadyExists) {
			return http.StatusConflict, nil
		}
		h.log.Log("%v", err)
		return http.StatusInternalServerError, nil
	}

	pointJson, err := json.Marshal(req)
	if err != nil {
		h.log.Log("%v", err)
		return http.StatusInternalServerError, nil
	}

	return http.StatusCreated, pointJson
}

func (h *PickUpPointHandlers) ListHandler(req *http.Request) (int, []byte) {
	list, err := h.svc.ListPoints(req.Context())
	if err != nil {
		h.log.Log("%v", err)
		return http.StatusInternalServerError, nil
	}

	listJson, err := json.Marshal(list)
	if err != nil {
		h.log.Log("%v", err)
		return http.StatusInternalServerError, nil
	}

	return http.StatusOK, listJson
}

func (h *PickUpPointHandlers) GetHandler(req *http.Request) (int, []byte) {
	idStr, ok := mux.Vars(req)["id"]
	if !ok {
		return http.StatusBadRequest, nil
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return http.StatusBadRequest, nil
	}

	point, err := h.svc.GetPoint(req.Context(), id)
	if err != nil {
		if errors.Is(err, pickuppoint.ErrNoItemFound) {
			return http.StatusNotFound, nil
		}
		h.log.Log("%v", err)
		return http.StatusInternalServerError, nil
	}

	pointJson, err := json.Marshal(point)
	if err != nil {
		h.log.Log("%v", err)
		return http.StatusInternalServerError, nil
	}

	return http.StatusOK, pointJson
}

func (h *PickUpPointHandlers) UpdateHandler(httpReq *http.Request) (int, []byte) {
	idStr, ok := mux.Vars(httpReq)["id"]
	if !ok {
		return http.StatusBadRequest, nil
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return http.StatusBadRequest, nil
	}

	body, err := io.ReadAll(httpReq.Body)
	if err != nil {
		h.log.Log("%v", err)
		return http.StatusInternalServerError, nil
	}

	var req core.UpdatePointRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		return http.StatusBadRequest, nil
	}

	if req.Id != id {
		return http.StatusBadRequest, nil
	}

	err = h.svc.UpdatePoint(httpReq.Context(), req)
	if err != nil {
		if errors.Is(err, pickuppoint.ErrNoItemFound) {
			return http.StatusNotFound, nil
		}
		h.log.Log("%v", err)
		return http.StatusInternalServerError, nil
	}

	pointJson, err := json.Marshal(req)
	if err != nil {
		h.log.Log("%v", err)
		return http.StatusInternalServerError, nil
	}

	return http.StatusOK, pointJson
}

func (h *PickUpPointHandlers) DeleteHandler(req *http.Request) (int, []byte) {
	idStr, ok := mux.Vars(req)["id"]
	if !ok {
		return http.StatusBadRequest, nil
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return http.StatusBadRequest, nil
	}

	err = h.svc.DeletePoint(req.Context(), id)
	if err != nil {
		if errors.Is(err, pickuppoint.ErrNoItemFound) {
			return http.StatusNotFound, nil
		}
		h.log.Log("%v", err)
		return http.StatusInternalServerError, nil
	}

	return http.StatusNoContent, nil
}
