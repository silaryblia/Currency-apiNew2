package service

import (
	"context"
	"time"

	"go.uber.org/zap"
)

func StartRatesScheduler(
	ctx context.Context,
	log *zap.Logger,
	svc *CurrencyService,
	interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				log.Info("scheduled rates sync started")
				if err := svc.SyncRates(ctx); err != nil {
					log.Error("scheduled rates sync failed", zap.Error(err))
				} else {
					log.Info("scheduled rates sync finished")
				}

			case <-ctx.Done():
				log.Info("rates scheduler stopped")
				return
			}
		}
	}()
}
