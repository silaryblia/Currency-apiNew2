package main

import (
	"Currency-apiNew2/internal/currency/provider"
	"Currency-apiNew2/internal/currency/repository"
	"Currency-apiNew2/internal/currency/service"
	"Currency-apiNew2/internal/currency/transport/grpc"
	_ "Currency-apiNew2/internal/currency/transport/http"
	"Currency-apiNew2/pkg/logger"
	_ "context"
	"time"

	_ "go.uber.org/zap"
)

func main() {
	log := logger.New()

	// Базовый провайдер ЦБ РФ
	baseProvider := provider.NewCBRProvider()

	// Кеш на 24 часа
	cachedProvider := provider.NewCachedProvider(baseProvider, 24*time.Hour)

	repo := repository.NewCurrencyRepoInMemory(log) // или Postgres

	svc := service.NewCurrencyService(repo, cachedProvider)

	grpc.RunServer(svc)

	//srv := http.NewServer(log)
	//if err := srv.Run(); err != nil {
	//
	//	log.Fatal("currency init failed", zap.Error(err))
	//}
}
