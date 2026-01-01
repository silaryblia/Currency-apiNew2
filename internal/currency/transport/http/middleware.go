package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}

func recoveryMiddleware(logger *zap.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error("PANIC recovered", zap.Any("error", rec))
					WriteJSON(w, http.StatusInternalServerError, nil, "internal currency error")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
