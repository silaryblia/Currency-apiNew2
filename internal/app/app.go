package app

import (
	"Currency-apiNew2/internal/config"
	"Currency-apiNew2/internal/currency/domain"
	"Currency-apiNew2/internal/currency/gateway"
	"Currency-apiNew2/internal/currency/notification"
	"Currency-apiNew2/internal/currency/provider"
	"Currency-apiNew2/internal/currency/repository"
	"Currency-apiNew2/internal/currency/service"
	"time"

	"go.uber.org/zap"
)

type App struct {
	Logger  *zap.Logger
	Gateway gateway.CurrencyGateway
	Config  *config.Config
	//Service *service.CurrencyService
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
		domain.DefaultNotificationConfig(),
	)

	gw := gateway.NewCurrencyGateway(svc, log)

	return &App{
		Logger:  log,
		Gateway: gw,
		Config:  cfg,
		//	Service: svc,
	}
}
