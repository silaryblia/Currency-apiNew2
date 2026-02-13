package http

import (
	"Currency-apiNew2/internal/currency/service"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func NewRouter(svc *service.CurrencyService, logger *zap.Logger) *mux.Router {
	r := mux.NewRouter()

	// Apply middlewares
	r.Use(JSONMiddleware)
	r.Use(recoveryMiddleware(logger))

	h := NewHandler(svc, logger)

	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/currencies/sync", h.SyncRates).Methods("POST")
	api.HandleFunc("/currencies", h.GetAll).Methods("GET")
	return r
}
