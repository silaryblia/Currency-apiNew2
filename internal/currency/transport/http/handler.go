package http

import (
	"Currency-apiNew2/internal/currency/domain"
	"Currency-apiNew2/internal/currency/service"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Handler struct {
	service *service.CurrencyService
	logger  *zap.Logger
}

func NewHandler(svc *service.CurrencyService, logger *zap.Logger) *Handler {
	return &Handler{service: svc, logger: logger}
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	rates, err := h.service.GetAll(r.Context())
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, nil, err.Error())
		return
	}
	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"date":  time.Now().Format(time.RFC822),
		"rates": rates,
	}, "")
}

func (h *Handler) GetOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := strings.ToUpper(vars["code"])

	rate, err := h.service.GetOne(r.Context(), code)
	if err != nil {
		WriteJSON(w, http.StatusNotFound, nil, "Currency not found")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"code": code,
		"rate": rate,
		"date": time.Now().Format(time.RFC822),
	}, "")
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Code string  `json:"code"`
		Rate float64 `json:"rate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSON(w, http.StatusBadRequest, nil, "invalid JSON body")
		return
	}

	if err := domain.ValidateCurrency(req.Code, req.Rate); err != nil {
		WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	date := time.Now()

	if req.Code == "" || req.Rate == 0 {
		WriteJSON(w, http.StatusBadRequest, nil, "invalid parameters")
		return
	}

	if err := h.service.Create(r.Context(), req.Code, req.Rate, date); err != nil {
		WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"code": req.Code,
		"rate": req.Rate,
		"date": date,
	}, "")
}

func (h *Handler) UpdateOne(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Code string  `json:"code"`
		Rate float64 `json:"rate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSON(w, http.StatusBadRequest, nil, "invalid JSON")
		return
	}

	if req.Code == "" {
		WriteJSON(w, http.StatusBadRequest, nil, "currency code required")
		return
	}

	date := time.Now()

	if err := h.service.UpdateOne(r.Context(), req.Code, req.Rate, date); err != nil {
		WriteJSON(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"code": req.Code,
		"rate": req.Rate,
		"date": date,
	}, "")
}

func (h *Handler) UpdateAll(w http.ResponseWriter, r *http.Request) {
	if err := h.service.UpdateAll(r.Context()); err != nil {
		WriteJSON(w, http.StatusInternalServerError, nil, err.Error())
		return
	}

	rates, _ := h.service.GetAll(r.Context())
	WriteJSON(w, http.StatusOK, rates, "")
}

func (h *Handler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	h.service.DeleteAll(r.Context())

	WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Все курсы удалены",
	}, "")
}

func (h *Handler) SyncRates(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("SyncRates called")

	if err := h.service.SyncRates(r.Context()); err != nil {
		WriteJSON(w, http.StatusInternalServerError, nil, err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{
		"message": "rates updated",
	}, "")
}
