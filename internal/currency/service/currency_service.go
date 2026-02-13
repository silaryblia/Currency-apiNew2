package service

import (
	"Currency-apiNew2/internal/currency/domain"
	"Currency-apiNew2/internal/currency/retry"
	"context"
	"sync/atomic"
	"time"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type CurrencyService struct {
	repo          domain.CurrencyRepository
	provider      domain.RatesProvider
	logger        *zap.Logger
	notifications domain.NotificationService
	notifyCfg     domain.NotificationConfig
	ready         atomic.Bool
}

func NewCurrencyService(
	repo domain.CurrencyRepository,
	provider domain.RatesProvider,
	notifications domain.NotificationService,
	logger *zap.Logger,
	notifyCfg domain.NotificationConfig,
) *CurrencyService {
	return &CurrencyService{
		repo:          repo,
		provider:      provider,
		notifications: notifications,
		notifyCfg:     notifyCfg,
		logger:        logger,
	}
}

func (s *CurrencyService) GetAll(ctx context.Context) (map[domain.CurrencyCode]domain.Currency, error) {
	s.logger.Info("getting all currencies")
	start := time.Now()

	res, err := s.repo.GetAll(ctx)
	elapsed := time.Since(start)
	if elapsed > s.notifyCfg.SlowThreshold {
		s.notifications.SlowOperation(ctx, "CurrencyService.GetAll", elapsed)
	}
	return res, err
}

func (s *CurrencyService) SyncRates(ctx context.Context) error {

	rawRates, rateDate, err := s.provider.ForceRefresh(ctx)
	if err != nil {
		return err
	}

	threshold := decimal.NewFromFloat(s.notifyCfg.RateSpikeThreshold)

	for rawCode, rawRate := range rawRates {
		code, err := domain.ParseCurrencyCode(rawCode)
		if err != nil {
			s.logger.Warn("invalid currency code", zap.String("code", rawCode), zap.Error(err))
			continue
		}

		newRate, err := domain.NewRate(rawRate)
		if err != nil {
			s.logger.Warn("invalid rate", zap.String("code", rawCode), zap.Float64("rate", rawRate), zap.Error(err))
			continue
		}

		old, err := s.repo.GetLatest(ctx, code)
		if err == nil {
			diff := newRate.Diff(old.Rate)
			if diff.GreaterThanOrEqual(threshold) {
				s.notifications.RateSpike(ctx, code, old.Rate, newRate)
			}
		}
		if err := s.repo.SaveRate(ctx, code, newRate, rateDate); err != nil {
			s.logger.Error("save rate failed",
				zap.String("code", rawCode),
				zap.Error(err))
		}
	}
	return nil
}

func (s *CurrencyService) GetLatest(
	ctx context.Context,
	code domain.CurrencyCode,
) (domain.Currency, error) {
	return s.repo.GetLatest(ctx, code)
}

func (s *CurrencyService) SetReady() {
	s.ready.Store(true)
}

func (s *CurrencyService) IsReady() bool {
	return s.ready.Load()
}

func (s *CurrencyService) SyncRatesWithRetry(ctx context.Context) error {
	return retry.WithBackoff(
		ctx,
		5,
		1*time.Second,
		1*time.Minute,
		func() error {
			return s.SyncRates(ctx)
		})
}

func (s *CurrencyService) GetAtDate(
	ctx context.Context,
	code domain.CurrencyCode,
	date time.Time,
) (domain.Currency, error) {
	return s.repo.GetAtDate(ctx, code, date)
}

func (s *CurrencyService) GetRange(
	ctx context.Context,
	code domain.CurrencyCode,
	from, to time.Time,
) ([]domain.Currency, error) {
	return s.repo.GetRange(ctx, code, from, to)
}
