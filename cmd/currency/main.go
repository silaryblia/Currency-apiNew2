package main

import (
	"Currency-apiNew2/internal/config"
	notification2 "Currency-apiNew2/internal/currency/notification"
	"Currency-apiNew2/internal/currency/provider"
	"Currency-apiNew2/internal/currency/repository"
	"Currency-apiNew2/internal/currency/service"
	"Currency-apiNew2/internal/currency/transport/grpc"
	"os"

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
		fmt.Fprintf(os.Stderr, "config load failed: %v\n", err)
		os.Exit(1)
	}

	log, err := logger.New(cfg.LogMode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logger init failed: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	log.Info("config loaded successfully")

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
