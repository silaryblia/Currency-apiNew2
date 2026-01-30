package service

import (
	"Currency-apiNew2/internal/currency/domain"
	"context"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
)

type CurrencyService struct {
	repo          domain.CurrencyRepository
	provider      domain.RatesProvider
	logger        *zap.Logger
	notifications domain.NotificationService
	mu            sync.RWMutex

	rateSpikeThreshold float64
	slowThreshold      time.Duration
	ready              atomic.Bool
}

func NewCurrencyService(
	repo domain.CurrencyRepository,
	provider domain.RatesProvider,
	notifications domain.NotificationService,
	logger *zap.Logger,
) *CurrencyService {
	return &CurrencyService{
		repo:               repo,
		provider:           provider,
		notifications:      notifications,
		rateSpikeThreshold: 0.1,                    // Например, 10%
		slowThreshold:      100 * time.Millisecond, // 100ms
		logger:             logger,
	}
}

func (s *CurrencyService) GetAll(ctx context.Context) (map[domain.CurrencyCode]domain.Currency, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.logger.Info("getting all currencies")
	start := time.Now()

	res, err := s.repo.GetAll(ctx)
	elapsed := time.Since(start)
	if elapsed > s.slowThreshold {
		s.notifications.SlowOperation(ctx, "CurrencyService.GetAll", elapsed)
	}
	return res, err
}

func (s *CurrencyService) GetOne(ctx context.Context, code domain.CurrencyCode) (domain.Currency, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.repo.GetOne(ctx, code)
}

func (s *CurrencyService) Create(ctx context.Context, code domain.CurrencyCode, rate domain.Rate, date time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := domain.ValidateCreateCurrency(code, rate); err != nil {
		return err
	}

	return s.repo.Create(ctx, code, rate, date)
}

func (s *CurrencyService) UpdateOne(ctx context.Context, code domain.CurrencyCode, rate domain.Rate, date time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.repo.UpdateOne(ctx, code, rate, date)
}

func (s *CurrencyService) UpdateAll(ctx context.Context) error {
	return s.repo.UpdateAll(ctx)
}

func (s *CurrencyService) DeleteAll(ctx context.Context) error {
	return s.repo.DeleteAll(ctx)
}

func (s *CurrencyService) SyncRates(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	rawRates, rateDate, err := s.provider.ForceRefresh(ctx)
	if err != nil {
		return err
	}

	for rawCode, rawRate := range rawRates {
		code, err := domain.ParseCurrencyCode(rawCode)
		if err != nil {
			return err
		}

		newRate, err := domain.NewRate(rawRate)
		if err != nil {
			return err
		}

		old, err := s.repo.GetOne(ctx, code)
		if err == nil {
			diff := float64(newRate-old.Rate) / float64(old.Rate)
			if diff < 0 {
				diff = -diff
			}

			if diff >= s.rateSpikeThreshold {
				s.notifications.RateSpike(ctx, code, old.Rate, newRate)
			}
		}

		if err := s.repo.Upsert(ctx, code, newRate, rateDate); err != nil {
			return err
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

func RetryWithBackoff(
	ctx context.Context,
	maxAttempts int,
	initialDelay time.Duration,
	maxDelay time.Duration,
	fn func() error,
) error {
	delay := initialDelay

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}

		if attempt == maxAttempts {
			return err
		}

		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}

		delay *= 2
		if delay > maxDelay {
			delay = maxDelay
		}
	}
	return nil
}

func (s *CurrencyService) SyncRatesWithRetry(ctx context.Context) error {
	return RetryWithBackoff(
		ctx,
		5,
		1*time.Second,
		1*time.Minute,
		func() error {
			return s.SyncRates(ctx)
		})
}
