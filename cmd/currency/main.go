package main

import (
	"Currency-apiNew2/internal/config"
	notification2 "Currency-apiNew2/internal/currency/notification"
	"Currency-apiNew2/internal/currency/provider"
	"Currency-apiNew2/internal/currency/repository"
	"Currency-apiNew2/internal/currency/service"
	"Currency-apiNew2/internal/currency/transport/grpc"
	//_ "Currency-apiNew2/internal/currency/transport/http"
	"Currency-apiNew2/pkg/logger"
	_ "context"
	"fmt"
	"time"

	"go.uber.org/zap"
	_ "go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		//cfg := config.DefaultConfig()
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	log, err := logger.New(cfg.LogMode)
	if err != nil {
		panic(fmt.Sprintf("Failed to create logger: %v", err))
	}
	defer log.Sync()

	repo := repository.NewCurrencyRepoInMemory(log)

	// Базовый провайдер ЦБ РФ
	baseProvider := provider.NewCBRProvider(&cfg.CBR)

	// Кеш на 24 часа
	cachedProvider := provider.NewCachedProvider(baseProvider, 24*time.Hour)

	notifications := notification2.NewLoggerNotificationService(log)

	svc := service.NewCurrencyService(repo, cachedProvider, notifications)

	log.Info("starting gRPC server", zap.String("port", cfg.GRPCPort))

	fmt.Println("CFG GRPC PORT =", cfg.GRPCPort)

	if err := grpc.RunServer(svc, cfg.GRPCPort); err != nil {
		log.Fatal("gRPC server failed", zap.Error(err))
	}
}
