package httpserv

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"homework/internal/app/core"
	coremocks "homework/internal/app/core/mocks"
	logmocks "homework/internal/app/logger/mocks"
	"homework/internal/app/pickuppoint"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type PickUpPointHandlersTestSuite struct {
	suite.Suite
}

func (s *PickUpPointHandlersTestSuite) Test_CreateHandler() {
	tests := []struct {
		name     string
		reqBody  string
		coreReq  core.CreatePointRequest
		coreErr  error
		wantCode int
		wantBody []byte
	}{
		{
			name:    "ok",
			reqBody: "{\"id\":1,\"name\":\"Generic pick-up point\",\"address\":\"5, Test st., Moscow\",\"contact\":\"test@example.com\"}",
			coreReq: core.CreatePointRequest{
				Id:      1,
				Name:    "Generic pick-up point",
				Address: "5, Test st., Moscow",
				Contact: "test@example.com",
			},
			wantCode: http.StatusCreated,
			wantBody: []byte("{\"id\":1,\"name\":\"Generic pick-up point\",\"address\":\"5, Test st., Moscow\",\"contact\":\"test@example.com\"}"),
		},
		{
			name:     "invalid body",
			reqBody:  "lhlihiuhilnjiklhni",
			wantCode: http.StatusBadRequest,
		},
		{
			name:    "existing id",
			reqBody: "{\"id\":1,\"name\":\"Generic pick-up point\",\"address\":\"5, Test st., Moscow\",\"contact\":\"test@example.com\"}",
			coreReq: core.CreatePointRequest{
				Id:      1,
				Name:    "Generic pick-up point",
				Address: "5, Test st., Moscow",
				Contact: "test@example.com",
			},
			coreErr:  pickuppoint.ErrIdAlreadyExists,
			wantCode: http.StatusConflict,
		},
		{
			name:    "error",
			reqBody: "{\"id\":1,\"name\":\"Generic pick-up point\",\"address\":\"5, Test st., Moscow\",\"contact\":\"test@example.com\"}",
			coreReq: core.CreatePointRequest{
				Id:      1,
				Name:    "Generic pick-up point",
				Address: "5, Test st., Moscow",
				Contact: "test@example.com",
			},
			coreErr:  assert.AnError,
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctrl := gomock.NewController(s.T())
			svc := coremocks.NewMockPickUpPointCoreService(ctrl)
			log := logmocks.NewMockLogger(ctrl)
			h := &PickUpPointHandlers{
				svc: svc,
				log: log,
			}
			svc.EXPECT().CreatePoint(gomock.Any(), tt.coreReq).Return(tt.coreErr).AnyTimes()
			log.EXPECT().Log(gomock.Any(), gomock.Any()).AnyTimes()
			req := httptest.NewRequest(http.MethodPost, "/pickup-points", strings.NewReader(tt.reqBody))
			code, body := h.CreateHandler(req, nil)
			s.Equal(tt.wantCode, code)
			s.Equal(tt.wantBody, body)
		})
	}
}

func (s *PickUpPointHandlersTestSuite) Test_ListHandler() {
	tests := []struct {
		name     string
		points   []pickuppoint.PickUpPoint
		coreErr  error
		wantCode int
		wantBody []byte
	}{
		{
			name: "ok",
			points: []pickuppoint.PickUpPoint{
				{
					Id:      1,
					Name:    "Generic pick-up point",
					Address: "5, Test st., Moscow",
					Contact: "test@example.com",
				},
				{
					Id:      2,
					Name:    "Another pick-up point",
					Address: "19, Sample st., Moscow",
					Contact: "sample@example.com",
				},
			},
			wantCode: http.StatusOK,
			wantBody: []byte("[{\"id\":1,\"name\":\"Generic pick-up point\",\"address\":\"5, Test st., Moscow\",\"contact\":\"test@example.com\"},{\"id\":2,\"name\":\"Another pick-up point\",\"address\":\"19, Sample st., Moscow\",\"contact\":\"sample@example.com\"}]"),
		},
		{
			name:     "error",
			coreErr:  assert.AnError,
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctrl := gomock.NewController(s.T())
			svc := coremocks.NewMockPickUpPointCoreService(ctrl)
			log := logmocks.NewMockLogger(ctrl)
			h := &PickUpPointHandlers{
				svc: svc,
				log: log,
			}
			svc.EXPECT().ListPoints(gomock.Any()).Return(tt.points, tt.coreErr)
			log.EXPECT().Log(gomock.Any(), gomock.Any()).AnyTimes()
			req := httptest.NewRequest(http.MethodGet, "/pickup-points", nil)
			code, body := h.ListHandler(req, nil)
			s.Equal(tt.wantCode, code)
			s.Equal(tt.wantBody, body)
		})
	}
}

func (s *PickUpPointHandlersTestSuite) Test_GetHandler() {
	tests := []struct {
		name     string
		idStr    string
		id       uint64
		point    pickuppoint.PickUpPoint
		coreErr  error
		wantCode int
		wantBody []byte
	}{
		{
			name:  "ok",
			idStr: "1",
			id:    1,
			point: pickuppoint.PickUpPoint{
				Id:      1,
				Name:    "Generic pick-up point",
				Address: "5, Test st., Moscow",
				Contact: "test@example.com",
			},
			wantCode: http.StatusOK,
			wantBody: []byte("{\"id\":1,\"name\":\"Generic pick-up point\",\"address\":\"5, Test st., Moscow\",\"contact\":\"test@example.com\"}"),
		},
		{
			name:     "invalid id string",
			idStr:    "dfsdfsf",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "not found",
			idStr:    "1",
			id:       1,
			coreErr:  pickuppoint.ErrNoItemFound,
			wantCode: http.StatusNotFound,
		},
		{
			name:     "error",
			idStr:    "1",
			id:       1,
			coreErr:  assert.AnError,
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctrl := gomock.NewController(s.T())
			svc := coremocks.NewMockPickUpPointCoreService(ctrl)
			log := logmocks.NewMockLogger(ctrl)
			h := &PickUpPointHandlers{
				svc: svc,
				log: log,
			}
			svc.EXPECT().GetPoint(gomock.Any(), tt.id).Return(tt.point, tt.coreErr).AnyTimes()
			log.EXPECT().Log(gomock.Any(), gomock.Any()).AnyTimes()
			req := httptest.NewRequest(http.MethodGet, "/pickup-points", nil)
			code, body := h.GetHandler(req, map[string]string{"id": tt.idStr})
			s.Equal(tt.wantCode, code)
			s.Equal(tt.wantBody, body)
		})
	}
}

func (s *PickUpPointHandlersTestSuite) Test_UpdateHandler() {
	tests := []struct {
		name     string
		idStr    string
		reqBody  string
		coreReq  core.UpdatePointRequest
		coreErr  error
		wantCode int
		wantBody []byte
	}{
		{
			name:    "ok",
			idStr:   "1",
			reqBody: "{\"id\":1,\"name\":\"Generic pick-up point\",\"address\":\"5, Test st., Moscow\",\"contact\":\"test@example.com\"}",
			coreReq: core.UpdatePointRequest{
				Id:      1,
				Name:    "Generic pick-up point",
				Address: "5, Test st., Moscow",
				Contact: "test@example.com",
			},
			wantCode: http.StatusOK,
			wantBody: []byte("{\"id\":1,\"name\":\"Generic pick-up point\",\"address\":\"5, Test st., Moscow\",\"contact\":\"test@example.com\"}"),
		},
		{
			name:     "invalid id string",
			idStr:    "dfsdfsf",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "invalid body",
			idStr:    "1",
			reqBody:  "lhlihiuhilnjiklhni",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "id and body mismatch",
			idStr:    "2",
			reqBody:  "{\"id\":1,\"name\":\"Generic pick-up point\",\"address\":\"5, Test st., Moscow\",\"contact\":\"test@example.com\"}",
			wantCode: http.StatusBadRequest,
		},
		{
			name:    "not found",
			idStr:   "1",
			reqBody: "{\"id\":1,\"name\":\"Generic pick-up point\",\"address\":\"5, Test st., Moscow\",\"contact\":\"test@example.com\"}",
			coreReq: core.UpdatePointRequest{
				Id:      1,
				Name:    "Generic pick-up point",
				Address: "5, Test st., Moscow",
				Contact: "test@example.com",
			},
			coreErr:  pickuppoint.ErrNoItemFound,
			wantCode: http.StatusNotFound,
		},
		{
			name:    "error",
			idStr:   "1",
			reqBody: "{\"id\":1,\"name\":\"Generic pick-up point\",\"address\":\"5, Test st., Moscow\",\"contact\":\"test@example.com\"}",
			coreReq: core.UpdatePointRequest{
				Id:      1,
				Name:    "Generic pick-up point",
				Address: "5, Test st., Moscow",
				Contact: "test@example.com",
			},
			coreErr:  assert.AnError,
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctrl := gomock.NewController(s.T())
			svc := coremocks.NewMockPickUpPointCoreService(ctrl)
			log := logmocks.NewMockLogger(ctrl)
			h := &PickUpPointHandlers{
				svc: svc,
				log: log,
			}
			svc.EXPECT().UpdatePoint(gomock.Any(), tt.coreReq).Return(tt.coreErr).AnyTimes()
			log.EXPECT().Log(gomock.Any(), gomock.Any()).AnyTimes()
			req := httptest.NewRequest(http.MethodPut, "/pickup-points", strings.NewReader(tt.reqBody))
			code, body := h.UpdateHandler(req, map[string]string{"id": tt.idStr})
			s.Equal(tt.wantCode, code)
			s.Equal(tt.wantBody, body)
		})
	}
}

func (s *PickUpPointHandlersTestSuite) Test_DeleteHandler() {
	tests := []struct {
		name     string
		idStr    string
		id       uint64
		coreErr  error
		wantCode int
	}{
		{
			name:     "ok",
			idStr:    "1",
			id:       1,
			wantCode: http.StatusNoContent,
		},
		{
			name:     "invalid id string",
			idStr:    "dfsdfsf",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "not found",
			idStr:    "1",
			id:       1,
			coreErr:  pickuppoint.ErrNoItemFound,
			wantCode: http.StatusNotFound,
		},
		{
			name:     "error",
			idStr:    "1",
			id:       1,
			coreErr:  assert.AnError,
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctrl := gomock.NewController(s.T())
			svc := coremocks.NewMockPickUpPointCoreService(ctrl)
			log := logmocks.NewMockLogger(ctrl)
			h := &PickUpPointHandlers{
				svc: svc,
				log: log,
			}
			svc.EXPECT().DeletePoint(gomock.Any(), tt.id).Return(tt.coreErr).AnyTimes()
			log.EXPECT().Log(gomock.Any(), gomock.Any()).AnyTimes()
			req := httptest.NewRequest(http.MethodDelete, "/pickup-points", nil)
			code, _ := h.DeleteHandler(req, map[string]string{"id": tt.idStr})
			s.Equal(tt.wantCode, code)
		})
	}
}

func TestPickUpPointHandlers(t *testing.T) {
	suite.Run(t, new(PickUpPointHandlersTestSuite))
}
