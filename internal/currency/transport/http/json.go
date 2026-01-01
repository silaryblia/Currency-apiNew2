package http

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}, errMsg string) {
	resp := APIResponse{
		Success: errMsg == "",
		Data:    data,
		Error:   errMsg,
	}

	w.Header().Set("Content-Type", "application/json")

	if rw, ok := w.(interface{ Written() bool }); !ok || !rw.Written() {
		w.WriteHeader(status)
	}

	_ = json.NewEncoder(w).Encode(resp)
}
