package main

import (
	"Currency-apiNew2/internal/app"
	"Currency-apiNew2/internal/currency/transport/grpc"
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM)
	defer stop()

	app := app.BuildApp()
	defer app.Shutdown()

	if err := app.Init(ctx); err != nil {
		app.Logger.Fatal("init failed", zap.Error(err))
	}

	go app.RunScheduler(ctx)
	go func() {
		if err := grpc.RunServer(ctx, app.Service, app.Config.GRPCPort); err != nil {
			app.Logger.Error("gRPC server error", zap.Error(err))
		}
	}()
	//go grpc.RunServer(ctx, app.Service, app.Config.GRPCPort)

	<-ctx.Done()
	app.Logger.Info("API service stopped")
}
