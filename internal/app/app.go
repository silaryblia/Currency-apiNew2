package app

import (
	"Currency-apiNew2/internal/config"
	"Currency-apiNew2/internal/currency/domain"
	"Currency-apiNew2/internal/currency/gateway"
	"Currency-apiNew2/internal/currency/notification"
	"Currency-apiNew2/internal/currency/provider"
	"Currency-apiNew2/internal/currency/repository"
	"Currency-apiNew2/internal/currency/service"
	"Currency-apiNew2/internal/metrics"
	"database/sql"
	"time"

	"go.uber.org/zap"
)

type App struct {
	Logger  *zap.Logger
	Gateway gateway.CurrencyGateway
	Config  *config.Config
	DB      *sql.DB
	//Service *service.CurrencyService
}

func BuildApp() *App {

	cfg := MustLoadConfig()
	log := MustInitLogger(cfg.LogMode)

	db, err := sql.Open("postgres", cfg.DB.DSN)
	if err != nil {
		log.Fatal("failed to open db", zap.Error(err))
	}

	if err := db.Ping(); err != nil {
		log.Fatal("failed to ping db", zap.Error(err))
	}

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
	metrics.MustRegister()

	return &App{
		Logger:  log,
		Gateway: gw,
		Config:  cfg,
		DB:      db,
		//	Service: svc,
	}
}
