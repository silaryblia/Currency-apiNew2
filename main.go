package main

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// Get - получить данные
// Post - отправить данные
// Patch/Put - обновить данные
// Delete - удалить данные

// Create Read Update Delete

const (
	DefaultUSD = 80.00
	DefaultEUR = 85.00
	DefaultAED = 20.00
)

// API errors
var (
	ErrNotFound      = errors.New("currency not found")
	ErrAlreadyExists = errors.New("currency already exists")
	ErrInvalidArg    = errors.New("invalid argument")
)

type CurrencyRepoInMemory struct {
	mu     sync.RWMutex
	data   map[string]float64
	logger *zap.Logger
}

type CurrencyServer struct {
	repo   *CurrencyRepoInMemory
	router *mux.Router
	logger *zap.Logger
}

type CreateOrUpdateCurrencyRequest struct {
	Currency string  `json:"code"`
	Rate     float64 `json:"rate"`
}

type UpdateCurrencyRequest struct {
	Currency string  `json:"code"`
	Rate     float64 `json:"rate"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// JSON Writer
func WriteJSON(w http.ResponseWriter, status int, data interface{}, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := APIResponse{
		Success: errMsg == "",
		Data:    data,
		Error:   errMsg,
	}

	json.NewEncoder(w).Encode(resp)
}

// JSON Middleware
func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}

// ////////////////////////////////////////////////////
// REPOSITORY
// ////////////////////////////////////////////////////
// Create
func NewCurrencyRepoInMemory(logger *zap.Logger) *CurrencyRepoInMemory {
	return &CurrencyRepoInMemory{
		data: map[string]float64{
			"USD": DefaultUSD,
			"EUR": DefaultEUR,
			"AED": DefaultAED,
		},
		logger: logger,
	}
}

// Read
func (r *CurrencyRepoInMemory) Get(code string) (float64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	code = strings.ToUpper(strings.TrimSpace(code))
	r.logger.Debug("repo: Get called", zap.String("code", code))

	v, ok := r.data[code]
	if !ok {
		r.logger.Warn("repo: currency not found", zap.String("code", code))
		return 0, ErrNotFound
	}

	r.logger.Debug("repo: Get success", zap.String("code", code), zap.Float64("rate", v))
	return v, nil
}

func (r *CurrencyRepoInMemory) GetAll() (map[string]float64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	r.logger.Debug("repo: GetAll called", zap.Int("count", len(r.data)))

	copy := make(map[string]float64, len(r.data))
	for k, v := range r.data {
		copy[k] = v
	}

	r.logger.Debug("repo: GetAll success", zap.Int("count", len(r.data)))
	return copy, nil
}

// add new currency
func (r *CurrencyRepoInMemory) Create(code string, rate float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	code = strings.ToUpper(strings.TrimSpace(code))

	if _, ok := r.data[code]; ok {
		return ErrAlreadyExists
	}

	r.data[code] = rate
	r.logger.Debug("repo: Create called", zap.String("code", code))
	return nil
}

// update one currency
func (r *CurrencyRepoInMemory) UpdateOne(code string, rate float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	code = strings.ToUpper(strings.TrimSpace(code))

	if _, ok := r.data[code]; !ok {
		return ErrNotFound
	}

	r.data[code] = rate
	r.logger.Debug("repo: UpdateOne called", zap.String("code", code))
	return nil

}

// Update
func (r *CurrencyRepoInMemory) UpdateAll() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.logger.Debug("repo: UpdateAll called", zap.Int("count", len(r.data)))

	// clear old, add new
	for k, v := range r.data {
		change := rand.Float64()*10 - 5 // [-5, 5]
		newVal := v + change
		if newVal < 0 {
			newVal = 0
		}
		r.data[k] = newVal
	}

	r.logger.Debug("repo: UpdateAll success", zap.Int("count", len(r.data)))
}

// Delete
func (r *CurrencyRepoInMemory) DeleteAll() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.logger.Debug("repo: DeleteAll called", zap.Int("count", len(r.data)))

	// clear map
	r.data = make(map[string]float64)

	r.logger.Debug("repo: DeleteAll success", zap.Int("count", len(r.data)))
	return nil
}

// ////////////////////////////////////////////////////
// SERVER
// ////////////////////////////////////////////////////
func NewCurrencyServer(logger *zap.Logger) (*CurrencyServer, error) {
	logger.Debug("repo: NewCurrencyServer called")

	repo := NewCurrencyRepoInMemory(logger)

	s := &CurrencyServer{
		repo:   repo,
		logger: logger,
	}

	s.router = s.NewRouter() // handlers registration

	logger.Debug("repo: NewCurrencyServer success")
	return s, nil
}

func (s *CurrencyServer) NewRouter() *mux.Router {
	r := mux.NewRouter()

	// Apply middlewares
	r.Use(JSONMiddleware)
	r.Use(recoveryMiddleware(s.logger))

	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/currency", s.handleGetAll).Methods("GET")
	api.HandleFunc("/currency/{code}", s.handleGetOne).Methods("GET")
	api.HandleFunc("/currency/create", s.handleCreate).Methods("POST")
	api.HandleFunc("/currency/update", s.handleUpdateOne).Methods("PUT")

	api.HandleFunc("/update", s.handleUpdate).Methods("PUT")
	api.HandleFunc("/delete", s.handleDelete).Methods("DELETE")

	return r
}

// HANDLERS
func (s *CurrencyServer) handleGetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteJSON(w, http.StatusMethodNotAllowed, nil, "Method not allowed")
		return
	}

	s.logger.Info("Обработка HTTP GET /currency", zap.String("method", r.Method))

	rates, err := s.repo.GetAll()
	if err != nil {
		s.logger.Error("failed to get all currencies", zap.Error(err))
		WriteJSON(w, http.StatusInternalServerError, nil, err.Error())

		return
	}

	s.logger.Info("Отправка курса валют клиенту", zap.Int("count", len(rates)))
	WriteJSON(w, http.StatusOK, rates, "")
}

func (s *CurrencyServer) handleGetOne(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteJSON(w, http.StatusMethodNotAllowed, nil, "Method not allowed")
		return
	}

	vars := mux.Vars(r)
	//code := strings.TrimPrefix(r.URL.Path, "/currency/")
	code := strings.ToUpper(vars["code"])

	s.logger.Info("Обработка HTTP GET /currencies/{code}", zap.String("code", code))

	rate, err := s.repo.Get(code)
	if err != nil {
		s.logger.Warn("Currency not found", zap.String("code", code))

		WriteJSON(w, http.StatusNotFound, nil, "Currency not found")
		return
	}

	s.logger.Info("Отправка курса валюты клиенту", zap.String("code", code), zap.Float64("rate", rate))

	data := map[string]interface{}{"code": code, "rate": rate}
	WriteJSON(w, http.StatusOK, data, "")
}

func (s *CurrencyServer) handleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteJSON(w, http.StatusMethodNotAllowed, nil, "Method not allowed")
		return
	}

	s.logger.Info("Обработка HTTP POST /create", zap.String("method", r.Method))

	var req CreateOrUpdateCurrencyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSON(w, http.StatusBadRequest, nil, "invalid JSON body")
		return
	}

	if req.Currency == "" || req.Rate == 0.0 {
		WriteJSON(w, http.StatusBadRequest, nil, "invalid parameters")
		return
	}

	err := s.repo.Create(req.Currency, req.Rate)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"code": req.Currency,
		"rate": req.Rate,
	}, "")
}

func (s *CurrencyServer) handleUpdateOne(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		WriteJSON(w, http.StatusMethodNotAllowed, nil, "Method not allowed")
		return
	}

	var req CreateOrUpdateCurrencyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSON(w, http.StatusBadRequest, nil, "invalid JSON")
		return
	}

	if req.Currency == "" {
		WriteJSON(w, http.StatusBadRequest, nil, "currency code required")
		return
	}

	err := s.repo.UpdateOne(req.Currency, req.Rate)
	if err != nil {
		WriteJSON(w, http.StatusNotFound, nil, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"code": req.Currency,
		"rate": req.Rate,
	}, "")
}

func (s *CurrencyServer) handleUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		WriteJSON(w, http.StatusMethodNotAllowed, nil, "Method not allowed")
		return
	}

	s.logger.Info("Обработка HTTP PUT /currencies")

	s.repo.UpdateAll()

	newValues, err := s.repo.GetAll()
	if err != nil {
		s.logger.Error("failed to update all currencies", zap.Error(err))

		WriteJSON(w, http.StatusInternalServerError, nil, "failed to update all currencies")
		return
	}

	s.logger.Info("Обновление курсов валют", zap.Int("count", len(newValues)))

	WriteJSON(w, http.StatusOK, newValues, "")
}

func (s *CurrencyServer) handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		WriteJSON(w, http.StatusMethodNotAllowed, nil, "Method not allowed")
		return
	}

	s.logger.Info("Обработка HTTP DELETE")

	s.repo.DeleteAll()

	s.logger.Info("Курсы валют удалены")

	WriteJSON(w, http.StatusOK, map[string]string{"message": "Все курсы удалены"}, "")
}

// /////////////////////////////////
// MAIN + graceful shutdown
// /////////////////////////////////
func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		logger.Sugar().Fatalf("failed to initialize zap logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Logger initialized")

	//repo := NewCurrencyRepoInMemory(map[string]float64{
	//	"USD": 80.0,
	//	"EUR": 85.0,
	//	"AED": 20.0,
	//}, logger)

	srv, err := NewCurrencyServer(logger)
	if err != nil {
		logger.Fatal("currency init failed", zap.Error(err))
	}

	if err := srv.Run(); err != nil {
		logger.Fatal("currency stopped with error", zap.Error(err))
	}
}

func (s *CurrencyServer) Run() error {
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: s.router,
	}

	go func() {
		s.logger.Info("Server starting...", zap.String("addr", httpServer.Addr))

		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("Server run error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	s.logger.Info("Shutdown received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return httpServer.Shutdown(ctx)
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
