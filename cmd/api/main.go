package main

import (
	"Currency-apiNew2/internal/app"
	pb "Currency-apiNew2/internal/currency/proto"
	"Currency-apiNew2/internal/currency/transport/grpc"
	currencyhttp "Currency-apiNew2/internal/currency/transport/http"
	"context"
	nethttp "net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type server struct {
	pb.CurrencyServiceServer
}

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM)
	defer stop()

	// metrics - ИСПОЛЬЗУЕМ nethttp
	go func() {
		nethttp.Handle("/metrics", promhttp.Handler())
		if err := nethttp.ListenAndServe(":2112", nil); err != nil {
			panic(err)
		}
	}()

	app := app.BuildApp()
	defer app.Shutdown()

	// ИСПОЛЬЗУЕМ currencyhttp ДЛЯ ТВОЕГО ХЕНДЛЕРА
	healthHandler := currencyhttp.NewHealthHandler(app)

	// ИСПОЛЬЗУЕМ nethttp ДЛЯ ВСЕГО СТАНДАРТНОГО
	mux := nethttp.NewServeMux()
	mux.HandleFunc("/health", healthHandler.Health)
	mux.HandleFunc("/ready", healthHandler.Ready)

	go func() {
		if err := nethttp.ListenAndServe(":8082", mux); err != nil {
			app.Logger.Error("health server failed", zap.Error(err))
		}
	}()

	if err := app.Init(ctx); err != nil {
		app.Logger.Fatal("init failed", zap.Error(err))
	}

	go app.RunScheduler(ctx)
	go func() {
		if err := grpc.RunServer(ctx, app.Gateway, app.Config.GRPCPort); err != nil {
			app.Logger.Error("gRPC server error", zap.Error(err))
		}
	}()

	<-ctx.Done()
	app.Logger.Info("API service stopped")
}
