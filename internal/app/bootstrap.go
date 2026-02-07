package app

import (
	"Currency-apiNew2/internal/config"
	"Currency-apiNew2/pkg/logger"
	"context"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
)

func MustLoadConfig() *config.Config {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config load failed: %v\n", err)
		os.Exit(1)
	}
	return &cfg
}

func MustInitLogger(mode string) *zap.Logger {
	log, err := logger.New(mode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logger init failed: %v\n", err)
		os.Exit(1)
	}
	return log
}

func (a *App) Init(ctx context.Context) error {

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := a.Gateway.SyncRates(ctx); err != nil {
		a.Logger.Error("initial rates sync failed", zap.Error(err))
	}

	return nil
}

func (a *App) RunScheduler(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				_ = a.Gateway.SyncRates(ctx)
			}
		}
	}()
}

func (a *App) Shutdown() {
	_ = a.Logger.Sync()
}
