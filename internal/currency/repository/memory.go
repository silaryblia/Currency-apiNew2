package repository

import (
	"Currency-apiNew2/internal/currency/domain"
	"context"
	_ "math/rand"
	"sync"
	"time"

	"go.uber.org/zap"
)

const (
	DefaultUSD = 80.00
	DefaultEUR = 85.00
	DefaultAED = 20.00
)

type CurrencyRepoInMemory struct {
	mu      sync.RWMutex
	latest  map[domain.CurrencyCode]domain.Currency
	history map[domain.CurrencyCode][]domain.Currency
	logger  *zap.Logger
}

func NewCurrencyRepoInMemory(logger *zap.Logger) *CurrencyRepoInMemory {
	return &CurrencyRepoInMemory{
		latest:  make(map[domain.CurrencyCode]domain.Currency),
		history: make(map[domain.CurrencyCode][]domain.Currency),
		logger:  logger,
	}
}

func (r *CurrencyRepoInMemory) GetAll(ctx context.Context) (map[domain.CurrencyCode]domain.Currency, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	res := make(map[domain.CurrencyCode]domain.Currency, len(r.latest))
	for k, v := range r.latest {
		res[k] = v
	}
	return res, nil
}

func (r *CurrencyRepoInMemory) SaveRate(
	ctx context.Context,
	code domain.CurrencyCode,
	rate domain.Rate,
	date time.Time) error {

	r.mu.Lock()
	defer r.mu.Unlock()

	c := domain.Currency{
		Code:     code,
		Rate:     rate,
		RateDate: date,
	}

	r.history[code] = append(r.history[code], c)
	r.latest[code] = c

	return nil
}

func (r *CurrencyRepoInMemory) GetLatest(
	ctx context.Context,
	code domain.CurrencyCode,
) (domain.Currency, error) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	c, ok := r.latest[code]
	if !ok {
		return domain.Currency{}, domain.ErrNotFound
	}
	return c, nil
}

func (r *CurrencyRepoInMemory) GetAtDate(
	ctx context.Context,
	code domain.CurrencyCode,
	date time.Time) (domain.Currency, error) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	h := r.history[code]
	for _, c := range h {
		//if c.RateDate.Equal(date) {
		//	return c, nil
		//}
		if sameDay(c.RateDate, date) {
			return c, nil
		}
	}
	return domain.Currency{}, domain.ErrNotFound
}

func sameDay(a, b time.Time) bool {
	y1, m1, d1 := a.Date()
	y2, m2, d2 := b.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func (r *CurrencyRepoInMemory) GetRange(
	ctx context.Context,
	code domain.CurrencyCode,
	from, to time.Time,
) ([]domain.Currency, error) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	var res []domain.Currency
	for _, c := range r.history[code] {
		if !c.RateDate.Before(from) && !c.RateDate.After(to) {
			res = append(res, c)
		}
	}
	if len(res) == 0 {
		return nil, domain.ErrNotFound
	}
	return res, nil
}
