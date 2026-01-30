package main

import (
	"Currency-apiNew2/internal/config"
	"Currency-apiNew2/internal/currency/notification"
	"Currency-apiNew2/internal/currency/provider"
	"Currency-apiNew2/internal/currency/repository"
	"Currency-apiNew2/internal/currency/service"
	"Currency-apiNew2/internal/currency/transport/grpc"
	"Currency-apiNew2/pkg/logger"
	"context"
	"fmt"
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
	defer func() {
		_ = log.Sync()
	}()

	log.Info("service starting")

	repo := repository.NewCurrencyRepoInMemory(log)
	baseProvider := provider.NewCBRProvider(&cfg.CBR)
	cachedProvider := provider.NewCachedProvider(baseProvider, 24*time.Hour)
	notifications := notification.NewLoggerNotificationService(log)
	svc := service.NewCurrencyService(repo, cachedProvider, notifications, log)

	startCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := svc.SyncRatesWithRetry(startCtx); err != nil {
		log.Error("initial sync rates failed", zap.Error(err))
	} else {
		log.Info("initial sync rates completed")
		svc.SetReady()
	}

	go service.StartRatesScheduler(ctx, log, svc, 24*time.Hour)

	go func() {
		log.Info("starting gRPC server", zap.String("port", cfg.GRPCPort))
		if err := grpc.RunServer(ctx, svc, cfg.GRPCPort); err != nil {
			log.Error("gRPC server stopped", zap.Error(err))
		}
	}()

	<-ctx.Done()
	log.Info("shutdown signal received, stopping...")
}
