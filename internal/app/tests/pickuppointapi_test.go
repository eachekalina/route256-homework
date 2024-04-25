//go:build integration

package tests

import (
	"context"
	"github.com/stretchr/testify/suite"
	"homework/cmd/app/httpserv"
	"homework/internal/app/core"
	"homework/internal/app/db"
	"homework/internal/app/logger"
	"homework/internal/app/pickuppoint"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type PickUpPointApiIntegrationTestSuite struct {
	suite.Suite
	db      *db.PostgresDatabase
	stopLog context.CancelFunc
	h       *httpserv.PickUpPointHandlers
}

func (s *PickUpPointApiIntegrationTestSuite) SetupSuite() {
	var err error
	tm, err := db.NewTransactionManager(context.Background())
	if err != nil {
		panic(err)
	}
	s.db = db.NewDatabase(tm)
	repo := pickuppoint.NewPostgresRepository(s.db)
	svc := pickuppoint.NewService(repo, tm)
	log := logger.NewLogger()
	var ctx context.Context
	ctx, s.stopLog = context.WithCancel(context.Background())
	go log.Run(ctx)
	coreSvc := core.NewPickUpPointCoreService(svc, log)
	s.h = httpserv.NewPickUpPointHandlers(coreSvc, log)
}

func (s *PickUpPointApiIntegrationTestSuite) TearDownSuite() {
	s.stopLog()
}

func (s *PickUpPointApiIntegrationTestSuite) SetupTest() {
	_, err := s.db.Exec(context.Background(), `INSERT INTO pickup_points VALUES
                              (1, 'Generic pick-up point', '5, Test st., Moscow', 'test@example.com'),
                              (2, 'Another pick-up point', '19, Sample st., Moscow', 'sample@example.com');`)
	if err != nil {
		panic(err)
	}
}

func (s *PickUpPointApiIntegrationTestSuite) TearDownTest() {
	_, err := s.db.Exec(context.Background(), "DELETE FROM pickup_points;")
	if err != nil {
		panic(err)
	}
}

func (s *PickUpPointApiIntegrationTestSuite) Test_CreateHandler() {
	tests := []struct {
		name     string
		reqBody  string
		wantCode int
		wantBody []byte
	}{
		{
			name:     "ok",
			reqBody:  "{\"id\":3,\"name\":\"Generic pick-up point\",\"address\":\"5, Test st., Moscow\",\"contact\":\"test@example.com\"}",
			wantCode: http.StatusCreated,
			wantBody: []byte("{\"id\":3,\"name\":\"Generic pick-up point\",\"address\":\"5, Test st., Moscow\",\"contact\":\"test@example.com\"}"),
		},
		{
			name:     "invalid body",
			reqBody:  "lhlihiuhilnjiklhni",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "existing id",
			reqBody:  "{\"id\":1,\"name\":\"Generic pick-up point\",\"address\":\"5, Test st., Moscow\",\"contact\":\"test@example.com\"}",
			wantCode: http.StatusConflict,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			req := httptest.NewRequest(http.MethodPost, "/pickup-points", strings.NewReader(tt.reqBody))
			code, body := s.h.CreateHandler(req, nil)
			s.Equal(tt.wantCode, code)
			s.Equal(tt.wantBody, body)
		})
	}
}

func (s *PickUpPointApiIntegrationTestSuite) Test_ListHandler() {
	tests := []struct {
		name     string
		wantCode int
		wantBody []byte
	}{
		{
			name:     "ok",
			wantCode: http.StatusOK,
			wantBody: []byte("[{\"id\":1,\"name\":\"Generic pick-up point\",\"address\":\"5, Test st., Moscow\",\"contact\":\"test@example.com\"},{\"id\":2,\"name\":\"Another pick-up point\",\"address\":\"19, Sample st., Moscow\",\"contact\":\"sample@example.com\"}]"),
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			req := httptest.NewRequest(http.MethodGet, "/pickup-points", nil)
			code, body := s.h.ListHandler(req, nil)
			s.Equal(tt.wantCode, code)
			s.Equal(tt.wantBody, body)
		})
	}
}

func (s *PickUpPointApiIntegrationTestSuite) Test_GetHandler() {
	tests := []struct {
		name     string
		idStr    string
		wantCode int
		wantBody []byte
	}{
		{
			name:     "ok",
			idStr:    "1",
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
			idStr:    "3",
			wantCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			req := httptest.NewRequest(http.MethodGet, "/pickup-points", nil)
			code, body := s.h.GetHandler(req, map[string]string{"id": tt.idStr})
			s.Equal(tt.wantCode, code)
			s.Equal(tt.wantBody, body)
		})
	}
}

func (s *PickUpPointApiIntegrationTestSuite) Test_UpdateHandler() {
	tests := []struct {
		name     string
		idStr    string
		reqBody  string
		wantCode int
		wantBody []byte
	}{
		{
			name:     "ok",
			idStr:    "1",
			reqBody:  "{\"id\":1,\"name\":\"Generic pick-up point updated\",\"address\":\"5, Updated st., Moscow\",\"contact\":\"updated@example.com\"}",
			wantCode: http.StatusOK,
			wantBody: []byte("{\"id\":1,\"name\":\"Generic pick-up point updated\",\"address\":\"5, Updated st., Moscow\",\"contact\":\"updated@example.com\"}"),
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
			reqBody:  "{\"id\":1,\"name\":\"Generic pick-up point updated\",\"address\":\"5, Updated st., Moscow\",\"contact\":\"updated@example.com\"}",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "not found",
			idStr:    "3",
			reqBody:  "{\"id\":3,\"name\":\"Generic pick-up point\",\"address\":\"5, Test st., Moscow\",\"contact\":\"test@example.com\"}",
			wantCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			req := httptest.NewRequest(http.MethodPut, "/pickup-points", strings.NewReader(tt.reqBody))
			code, body := s.h.UpdateHandler(req, map[string]string{"id": tt.idStr})
			s.Equal(tt.wantCode, code)
			s.Equal(tt.wantBody, body)
		})
	}
}

func (s *PickUpPointApiIntegrationTestSuite) Test_DeleteHandler() {
	tests := []struct {
		name     string
		idStr    string
		wantCode int
	}{
		{
			name:     "ok",
			idStr:    "1",
			wantCode: http.StatusNoContent,
		},
		{
			name:     "invalid id string",
			idStr:    "dfsdfsf",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "not found",
			idStr:    "3",
			wantCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			req := httptest.NewRequest(http.MethodDelete, "/pickup-points", nil)
			code, _ := s.h.DeleteHandler(req, map[string]string{"id": tt.idStr})
			s.Equal(tt.wantCode, code)
		})
	}
}

func TestPickUpPointApi(t *testing.T) {
	suite.Run(t, new(PickUpPointApiIntegrationTestSuite))
}
