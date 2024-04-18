package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"homework/internal/app/kafka"
	"homework/internal/app/logger"
	"homework/internal/app/middleware"
	"homework/internal/app/reqlog"
	"io"
	log2 "log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

var brokers = []string{
	"127.0.0.1:9191",
	"127.0.0.1:9192",
	"127.0.0.1:9193",
}

const topic = "requests"

func TestLogMiddleware(t *testing.T) {
	sarama.Logger = log2.New(os.Stdout, "[sarama]", log2.LstdFlags)

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	})

	resChan := make(chan reqlog.Message)

	logCtx, logStop := context.WithCancel(context.Background())
	defer logStop()
	log := logger.NewLogger()
	go log.Run(logCtx)
	prod, err := kafka.NewProducer(brokers, log, topic)
	assert.NoError(t, err)
	defer prod.Close()
	cons, err := kafka.NewConsumer(brokers, topic, func(bytes []byte) {
		var msg reqlog.Message
		_ = json.Unmarshal(bytes, &msg)

		resChan <- msg
	})
	assert.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		err := cons.Run(ctx)
		if err != nil {
			fmt.Println(err)
		}
	}()
	select {
	case <-time.After(15 * time.Second):
		assert.FailNow(t, "failed to init consumer", ctx.Err())
	case <-cons.Ready():
	}
	reqLog := reqlog.NewLogger(prod, cons)

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
			m := middleware.LogMiddleware(reqLog)
			h := m(h)
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			req = mux.SetURLVars(req, tt.params)
			req.Header = tt.headers
			w := httptest.NewRecorder()
			h.ServeHTTP(w, req)

			ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
			defer cancel()
			select {
			case <-ctx.Done():
				assert.Fail(t, "context cancelled", ctx.Err())
			case msg := <-resChan:
				assert.Equal(t, tt.method, msg.Method)
				assert.Equal(t, tt.path, msg.Path)
				assert.Equal(t, tt.headers, msg.Headers)
				assert.Equal(t, tt.params, msg.Params)
				assert.Equal(t, tt.body, msg.Body)
			}
		})
	}
}
