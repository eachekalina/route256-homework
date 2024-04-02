package middleware

import (
	"github.com/gorilla/mux"
	"net/http"
)

func AuthMiddleware(username string, password string) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			reqUsername, reqPassword, ok := req.BasicAuth()
			if !ok || reqUsername != username || reqPassword != password {
				w.Header().Set("WWW-Authenticate", "Basic realm=\"pickup-point\"")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			h.ServeHTTP(w, req)
		})
	}
}
