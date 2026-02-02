package service

import (
	"Currency-apiNew2/internal/currency/domain"
	"Currency-apiNew2/internal/currency/retry"
	"context"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
)

type CurrencyService struct {
	repo          domain.CurrencyRepository
	provider      domain.RatesProvider
	logger        *zap.Logger
	notifications domain.NotificationService
	//mu            sync.RWMutex

	notifyCfg domain.NotificationConfig
	ready     atomic.Bool
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
	if elapsed > s.notifyCfg.SlowThresold {
		s.notifications.SlowOperation(ctx, "CurrencyService.GetAll", elapsed)
	}
	return res, err
}

func (s *CurrencyService) GetOne(ctx context.Context, code domain.CurrencyCode) (domain.Currency, error) {
	return s.repo.GetOne(ctx, code)
}

func (s *CurrencyService) Create(ctx context.Context, code domain.CurrencyCode, rate domain.Rate, date time.Time) error {
	if err := domain.ValidateCreateCurrency(code, rate); err != nil {
		return err
	}

	return s.repo.Create(ctx, code, rate, date)
}

func (s *CurrencyService) UpdateOne(ctx context.Context, code domain.CurrencyCode, rate domain.Rate, date time.Time) error {
	return s.repo.UpdateOne(ctx, code, rate, date)
}

func (s *CurrencyService) UpdateAll(ctx context.Context) error {
	return s.repo.UpdateAll(ctx)
}

func (s *CurrencyService) DeleteAll(ctx context.Context) error {
	return s.repo.DeleteAll(ctx)
}

func (s *CurrencyService) SyncRates(ctx context.Context) error {
	rawRates, rateDate, err := s.provider.ForceRefresh(ctx)
	if err != nil {
		return err
	}

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

		old, err := s.repo.GetOne(ctx, code)
		if err == nil {
			oldRateFloat := old.Rate.Float64()
			if oldRateFloat > 0 {
				newRateFloat := newRate.Float64()
				diff := (newRateFloat - oldRateFloat) / oldRateFloat
				if diff < 0 {
					diff = -diff
				}
				if diff >= s.notifyCfg.RateSpikeThreshold {
					s.notifications.RateSpike(ctx, code, old.Rate, newRate)
				}
			}
		}

		if err := s.repo.Upsert(ctx, code, newRate, rateDate); err != nil {
			s.logger.Error("failed to upsert", zap.String("code", rawCode), zap.Error(err))
			continue
		}
	}
	return nil
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
