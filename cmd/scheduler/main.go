package main

import (
	"Currency-apiNew2/internal/app"
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
	//defer func() {
	//	_ = app.Logger.Sync()
	//}()

	app.Logger.Info("starting scheduler")

	if err := app.Gateway.SyncRates(ctx); err != nil {
		app.Logger.Error("init sync failed", zap.Error(err))
	} else {
		app.RunScheduler(ctx)
	}

	//service.StartRatesScheduler(
	//	ctx,
	//	app.Logger,
	//	app.Service,
	//	24*time.Hour)

	<-ctx.Done()
	app.Logger.Info("scheduler stopped")
}
