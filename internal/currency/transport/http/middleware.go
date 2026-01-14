package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type responseWriter struct {
	http.ResponseWriter
	written bool
}

func (w *responseWriter) WriteHeader(status int) {
	w.written = true
	w.ResponseWriter.WriteHeader(status)
}

func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func recoveryMiddleware(logger *zap.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := &responseWriter{ResponseWriter: w}

			defer func() {
				if rec := recover(); rec != nil {
					logger.Error("PANIC recovered", zap.Any("panic", rec))
					if !rw.written {
						_ = WriteError(w, http.StatusInternalServerError, "internal server error")
					}
				}
			}()
			next.ServeHTTP(rw, r)
		})
	}
}
