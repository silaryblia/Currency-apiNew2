package main

import (
	"Currency-apiNew2/internal/app"
	"Currency-apiNew2/internal/currency/transport/grpc"
	currencyhttp "Currency-apiNew2/internal/currency/transport/http"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	app := app.BuildApp()
	defer app.Shutdown()

	if err := app.Init(ctx); err != nil {
		app.Logger.Fatal("init failed", zap.Error(err))
	}

	mux := http.NewServeMux()

	healthHandler := currencyhttp.NewHealthHandler(app)
	mux.HandleFunc("/health", healthHandler.Health)
	mux.HandleFunc("/ready", healthHandler.Ready)
	mux.Handle("/metrics", promhttp.Handler())

	httpSrv := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	go func() {
		app.Logger.Info("HTTP server started on :8081")
		if err := httpSrv.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			app.Logger.Fatal("http server failed", zap.Error(err))
		}
	}()

	go app.RunScheduler(ctx)

	go func() {
		app.Logger.Info("gRPC server started",
			zap.String("port", app.Config.GRPCPort),
		)

		if err := grpc.RunServer(ctx, app.Gateway, app.Config.GRPCPort); err != nil {
			app.Logger.Error("gRPC server error", zap.Error(err))
		}
	}()

	<-ctx.Done()
	app.Logger.Info("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = httpSrv.Shutdown(shutdownCtx)

	app.Logger.Info("API stopped")
}
