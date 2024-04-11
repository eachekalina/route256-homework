package middleware

import (
	"bytes"
	"github.com/gorilla/mux"
	"homework/internal/app/logger"
	"io"
	"net/http"
)

func LogMiddleware(log logger.Logger) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			var buf bytes.Buffer
			req.Body = io.NopCloser(io.TeeReader(req.Body, &buf))
			h.ServeHTTP(w, req)
			log.Log("Request from %s, method %s, path %s, headers %s, body %s", req.RemoteAddr, req.Method, req.RequestURI, req.Header, buf.String())
		})
	}
}
