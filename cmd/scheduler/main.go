package main

import (
	"Currency-apiNew2/internal/app"
	"Currency-apiNew2/internal/currency/service"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM)
	defer stop()

	app := app.BuildApp()
	defer func() {
		_ = app.Logger.Sync()
	}()

	app.Logger.Info("starting scheduler")

	if err := app.Service.SyncRatesWithRetry(ctx); err != nil {
		app.Logger.Error("init sync failed", zap.Error(err))
	} else {
		app.Service.SetReady()
	}

	service.StartRatesScheduler(
		ctx,
		app.Logger,
		app.Service,
		24*time.Hour)

	<-ctx.Done()
	app.Logger.Info("scheduler stopped")
}
