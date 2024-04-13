package middleware

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"homework/internal/app/middleware/mocks"
	"homework/internal/app/reqlog"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogMiddleware(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name    string
		method  string
		path    string
		headers http.Header
		params  map[string]string
		body    string
	}{
		{
			name:    "basic",
			method:  http.MethodGet,
			path:    "/",
			headers: map[string][]string{},
		},
		{
			name:    "path",
			method:  http.MethodPost,
			path:    "/hello",
			headers: map[string][]string{},
		},
		{
			name:   "headers",
			method: http.MethodGet,
			path:   "/",
			headers: map[string][]string{
				"Authorization": {"Test"},
			},
		},
		{
			name:    "params",
			method:  http.MethodGet,
			path:    "/",
			headers: map[string][]string{},
			params: map[string]string{
				"id": "43",
			},
		},
		{
			name:    "body",
			method:  http.MethodGet,
			path:    "/",
			headers: map[string][]string{},
			body:    "somebodyuwu",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			log := mocks.NewMockRequestLogger(ctrl)
			log.EXPECT().Log(gomock.Any()).Do(func(msg reqlog.Message) {
				assert.Equal(t, tt.method, msg.Method)
				assert.Equal(t, tt.path, msg.Path)
				assert.Equal(t, tt.headers, msg.Headers)
				assert.Equal(t, tt.params, msg.Params)
				assert.Equal(t, tt.body, msg.Body)
			})
			m := LogMiddleware(log)
			h := m(h)
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			req = mux.SetURLVars(req, tt.params)
			req.Header = tt.headers
			w := httptest.NewRecorder()
			h.ServeHTTP(w, req)
		})
	}
}
