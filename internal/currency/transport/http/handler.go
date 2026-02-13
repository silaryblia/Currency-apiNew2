package http

import (
	"Currency-apiNew2/internal/currency/service"
	"net/http"
	"time"

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
		h.logger.Error("get all currencies failed", zap.Error(err))
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := map[string]any{
		"date":  time.Now().UTC().Format(time.RFC3339),
		"rates": rates,
	}

	if err := WriteJSON(w, http.StatusOK, resp); err != nil {
		h.logger.Error("failed to write response", zap.Error(err))
	}
}

func (h *Handler) SyncRates(w http.ResponseWriter, r *http.Request) {
	if err := h.service.SyncRates(r.Context()); err != nil {
		h.logger.Error("failed to sync currency rates", zap.Error(err))
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = WriteJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}
