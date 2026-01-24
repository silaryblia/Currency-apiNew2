package http

import (
	"Currency-apiNew2/internal/currency/domain"
	"Currency-apiNew2/internal/currency/service"
	"encoding/json"
	"net/http"
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

func (h *Handler) GetOne(w http.ResponseWriter, r *http.Request) {
	raw := mux.Vars(r)["code"]

	code, err := domain.ParseCurrencyCode(raw)
	if err != nil {
		_ = WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	cur, err := h.service.GetOne(r.Context(), code)
	if err != nil {
		if err == domain.ErrNotFound {
			_ = WriteError(w, http.StatusNotFound, "currency not found")
			return
		}

		h.logger.Error("get currency failed", zap.Error(err))
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := WriteJSON(w, http.StatusOK, cur); err != nil {
		h.logger.Error("write response failed", zap.Error(err))
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Code string  `json:"code"`
		Rate float64 `json:"rate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	code, err := domain.ParseCurrencyCode(req.Code)
	if err != nil {
		_ = WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	rate, err := domain.NewRate(req.Rate)
	if err != nil {
		_ = WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Create(r.Context(), code, rate, time.Now().UTC()); err != nil {
		status := HTTPStatusFromError(err)

		if status == http.StatusInternalServerError {
			h.logger.Error("create currency failed", zap.Error(err))
		}

		_ = WriteError(w, status, err.Error())
		return

		_ = WriteJSON(w, http.StatusCreated, map[string]string{
			"status": "ok",
		})
		//
		//	if err == domain.ErrAlreadyExists {
		//		_ = WriteError(w, http.StatusConflict, "currency already exists")
		//		return
		//	}
		//
		//	h.logger.Error("create currency failed", zap.Error(err))
		//	_ = WriteError(w, http.StatusInternalServerError, err.Error())
		//	return
		//}
		//
		//_ = WriteJSON(w, http.StatusCreated, map[string]string{
		//	"status": "ok",
		//})
	}
}

func (h *Handler) UpdateOne(w http.ResponseWriter, r *http.Request) {
	raw := mux.Vars(r)["code"]

	code, err := domain.ParseCurrencyCode(raw)
	if err != nil {
		_ = WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req struct {
		Rate float64 `json:"rate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	rate, err := domain.NewRate(req.Rate)
	if err != nil {
		_ = WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.UpdateOne(r.Context(), code, rate, time.Now().UTC()); err != nil {
		if err == domain.ErrNotFound {
			_ = WriteError(w, http.StatusNotFound, "currency not found")
			return
		}

		h.logger.Error("update currency failed", zap.Error(err))
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = WriteJSON(w, http.StatusOK, map[string]any{
		"code": code.String(),
		"rate": rate,
	})
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

func (h *Handler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	if err := h.service.DeleteAll(r.Context()); err != nil {
		h.logger.Error("failed to delete currency rates", zap.Error(err))
		_ = WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = WriteJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}
