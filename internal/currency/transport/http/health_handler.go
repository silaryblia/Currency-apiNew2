package http

import (
	"Currency-apiNew2/internal/app"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type HealthHandler struct {
	app *app.App
}

func NewHealthHandler(a *app.App) *HealthHandler {
	return &HealthHandler{app: a}
}

type healthResponse struct {
	Status string `json:"status"`
	DB     string `json:"db"`
	Ready  bool   `json:"ready"`
	Now    string `json:"now"`
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := r.Context(), func() {}
	if deadline, ok := ctx.Deadline(); !ok || time.Until(deadline) > 2*time.Second {
		ctx, cancel = context.WithTimeout(r.Context(), 2*time.Second)
	}
	defer cancel()

	resp := healthResponse{
		Status: "ok",
		DB:     "up",
		Ready:  h.app.Gateway.IsReady(),
		Now:    time.Now().UTC().Format(time.RFC3339),
	}

	if err := h.app.DB.PingContext(ctx); err != nil {
		resp.Status = "degraded"
		resp.DB = "down"
		w.WriteHeader(http.StatusServiceUnavailable)
	} else if !resp.Ready {
		resp.Status = "starting"
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	if !h.app.Gateway.IsReady() {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
}
