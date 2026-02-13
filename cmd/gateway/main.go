package main

import (
	"Currency-apiNew2/internal/currency/gateway"
	pb "Currency-apiNew2/internal/currency/proto"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(
		"currency-api:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("grpc close error:", err)
		}
	}()

	client := pb.NewCurrencyServiceClient(conn)
	handler := gateway.NewHandler(client)

	mux := http.NewServeMux()

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			"http://auth-service:8081/login",
			bytes.NewBuffer(body),
		)
		if err != nil {
			http.Error(w, "auth service error", http.StatusInternalServerError)
			return
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, "auth service unavailable", http.StatusInternalServerError)
			return
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Println("body close error:", err)
			}
		}()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		if _, err := io.Copy(w, resp.Body); err != nil {
			log.Println("copy error:", err)
		}
	})

	mux.HandleFunc(
		"/currencies",
		gateway.JWTMiddleware(handler.GetCurrencies),
	)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		fmt.Println("gateway started on :8080") // ← ТЕПЕРЬ ПРАВДА!
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = srv.Shutdown(shutdownCtx)
	fmt.Println("gateway stopped")
}
