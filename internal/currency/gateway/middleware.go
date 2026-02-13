package gateway

import (
	"context"
	"net/http"
	"time"
)

func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			"http://auth-service:8081/validate",
			nil,
		)
		if err != nil {
			http.Error(w, "auth service error", http.StatusInternalServerError)
			return
		}

		req.Header.Set("Authorization", authHeader)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, "auth service unavailable", http.StatusInternalServerError)
			return
		}
		defer func() {
			_ = resp.Body.Close()
		}()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
