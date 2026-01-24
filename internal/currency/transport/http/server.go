package http

import (
	"Currency-apiNew2/internal/config"
	_ "Currency-apiNew2/internal/config"
	"Currency-apiNew2/internal/currency/domain"
	"Currency-apiNew2/internal/currency/notification"
	"Currency-apiNew2/internal/currency/provider"
	"Currency-apiNew2/internal/currency/repository"
	"Currency-apiNew2/internal/currency/service"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type Server struct {
	router *mux.Router
	logger *zap.Logger
}

func NewServer(cfg *config.Config, logger *zap.Logger) (*Server, error) {
	var repo domain.CurrencyRepository

	if cfg.UsePostgres {
		db, err := sql.Open("postgres", cfg.PostgresDSN)
		if err != nil {
			return nil, fmt.Errorf("open postgres: %w", err)
		}
		if err := db.Ping(); err != nil {
			return nil, fmt.Errorf("ping postgres: %w", err)
		}

		repo = repository.NewCurrencyRepoPostgres(db, logger)
	} else {
		repo = repository.NewCurrencyRepoInMemory(logger)
	}

	// Базовый провайдер ЦБ РФ
	baseProvider := provider.NewCBRProvider(&cfg.CBR)
	// Кеш на 24 часа
	cachedProvider := provider.NewCachedProvider(baseProvider, 24*time.Hour)
	notificationSvc := notification.NewLoggerNotificationService(logger)

	// В сервис передаём КЕШ
	svc := service.NewCurrencyService(repo, cachedProvider, notificationSvc)
	r := NewRouter(svc, logger)

	return &Server{router: r, logger: logger}, nil
}

func (s *Server) Run() error {
	httpServer := &http.Server{
		Addr:    ":50051",
		Handler: s.router,
	}

	go func() {
		s.logger.Info("Server starting...", zap.String("addr", httpServer.Addr))

		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("Server run error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return httpServer.Shutdown(ctx)
}
