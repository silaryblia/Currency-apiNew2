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
	api.HandleFunc("/currencies/{code}", h.GetOne).Methods("GET")
	api.HandleFunc("/currencies/{code}", h.Create).Methods("POST")
	api.HandleFunc("/currencies/{code}", h.UpdateOne).Methods("PUT")
	api.HandleFunc("/currencies", h.UpdateAll).Methods("PUT")
	api.HandleFunc("/currencies", h.DeleteAll).Methods("DELETE")

	return r
}

//func NewRouter(logger *zap.Logger) *mux.Router {
//	r := mux.NewRouter()
//	r.Use(JSONMiddleware)
//	r.Use(recoveryMiddleware(logger))
//
//	var repo domain.CurrencyRepository
//
//	if os.Getenv("USE_POSTGRES") == "true" {
//		db, err := sql.Open("postgres", os.Getenv("POSTGRES_DSN"))
//		if err != nil {
//			logger.Fatal("postgres init failed", zap.Error(err))
//		}
//		repo = repository.NewCurrencyRepoPostgres(db, logger)
//	} else {
//		repo = repository.NewCurrencyRepoInMemory(logger)
//	}
//
//	svc := service.NewCurrencyService(repo)
//	h := NewHandler(svc, logger)
//
//	api := r.PathPrefix("/api/v1").Subrouter()
//	api.HandleFunc("/currency", h.GetAll).Methods("GET")
//	api.HandleFunc("/currency/{code}", h.GetOne).Methods("GET")
//	api.HandleFunc("/currency/create", h.Create).Methods("POST")
//	api.HandleFunc("/currency/update", h.UpdateOne).Methods("PUT")
//	api.HandleFunc("/update", h.UpdateAll).Methods("PUT")
//	api.HandleFunc("/delete", h.DeleteAll).Methods("DELETE")
//
//	return r
//}
