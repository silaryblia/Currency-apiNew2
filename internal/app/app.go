package app

import (
	"Currency-apiNew2/internal/config"
	"Currency-apiNew2/internal/currency/domain"
	"Currency-apiNew2/internal/currency/notification"
	"Currency-apiNew2/internal/currency/provider"
	"Currency-apiNew2/internal/currency/repository"
	"Currency-apiNew2/internal/currency/service"
	"time"

	"go.uber.org/zap"
)

type App struct {
	Logger  *zap.Logger
	Service *service.CurrencyService
	Config  *config.Config
}

func BuildApp() *App {
	cfg := MustLoadConfig()
	log := MustInitLogger(cfg.LogMode)
	repo := repository.NewCurrencyRepoInMemory(log)

	baseProvider := provider.NewCBRProvider(&cfg.CBR)
	cachedProvider := provider.NewCachedProvider(baseProvider, 1*time.Hour)

	notifications := notification.NewLoggerNotificationService(log)

	svc := service.NewCurrencyService(
		repo,
		cachedProvider,
		notifications,
		log,
		domain.NotificationConfig{
			RateSpikeThreshold: 0.1,
			SlowThresold:       100 * time.Millisecond,
		},
	)

	return &App{
		Logger:  log,
		Service: svc,
		Config:  cfg,
	}
}
