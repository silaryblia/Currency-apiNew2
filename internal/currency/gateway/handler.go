package gateway

import (
	pb "Currency-apiNew2/internal/currency/proto"
	"encoding/json"
	"net/http"
)

type Handler struct {
	Client pb.CurrencyServiceClient
}

func NewHandler(client pb.CurrencyServiceClient) *Handler {
	return &Handler{Client: client}
}

func (h *Handler) GetCurrencies(w http.ResponseWriter, r *http.Request) {
	resp, err := h.Client.GetAll(r.Context(), &pb.GetAllRequest{})
	if err != nil {
		http.Error(w, "currency error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	currencies := make([]*pb.Currency, 0, len(resp.Currencies))
	for _, curr := range resp.Currencies {
		currencies = append(currencies, curr)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(currencies); err != nil {
		http.Error(w, "encoding error", http.StatusInternalServerError)
		return
	}
}
