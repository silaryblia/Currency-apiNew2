package service

import (
	"Currency-apiNew2/internal/currency/domain"
	"context"
	"time"

	"go.uber.org/zap"
)

type CurrencyService struct {
	repo     domain.CurrencyRepository
	provider domain.RatesProvider
	logger   *zap.Logger
}

func NewCurrencyService(repo domain.CurrencyRepository, provider domain.RatesProvider) *CurrencyService {
	return &CurrencyService{
		repo:     repo,
		provider: provider,
	}
}

func (s *CurrencyService) GetAll(ctx context.Context) (map[domain.CurrencyCode]domain.Currency, error) {
	return s.repo.GetAll(ctx)
}

func (s *CurrencyService) GetOne(ctx context.Context, code domain.CurrencyCode) (domain.Currency, error) {
	return s.repo.GetOne(ctx, code)
}

func (s *CurrencyService) Create(ctx context.Context, code domain.CurrencyCode, rate domain.Rate, date time.Time) error {
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
			return err
		}

		rate, err := domain.NewRate(rawRate)
		if err != nil {
			return err
		}

		if err := s.repo.Upsert(ctx, code, rate, rateDate); err != nil {
			return err
		}
	}
	return nil
}
