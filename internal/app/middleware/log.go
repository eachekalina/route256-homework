//go:generate mockgen -source=./log.go -destination=./mocks/log.go -package=mocks

package middleware

import (
	"bytes"
	"github.com/gorilla/mux"
	"homework/internal/app/reqlog"
	"io"
	"net/http"
	"time"
)

type RequestLogger interface {
	Log(msg reqlog.Message)
}

func LogMiddleware(log RequestLogger) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			var buf bytes.Buffer
			req.Body = io.NopCloser(io.TeeReader(req.Body, &buf))
			h.ServeHTTP(w, req)

			msg := reqlog.Message{
				Timestamp: time.Now(),
				Method:    req.Method,
				Path:      req.RequestURI,
				Headers:   req.Header,
				Params:    mux.Vars(req),
				Body:      buf.String(),
			}
			log.Log(msg)
		})
	}
}
