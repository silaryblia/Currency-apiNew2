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

func (r *CurrencyRepoInMemory) GetOne(ctx context.Context,
	code domain.CurrencyCode) (domain.Currency, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	c, ok := r.latest[code]
	if !ok {
		return domain.Currency{}, domain.ErrNotFound
	}
	return c, nil
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

func (r *CurrencyRepoInMemory) Upsert(
	ctx context.Context,
	code domain.CurrencyCode,
	rate domain.Rate,
	date time.Time) error {
	return r.SaveRate(ctx, code, rate, date)
}

// add new currency
func (r *CurrencyRepoInMemory) Create(ctx context.Context, code domain.CurrencyCode, rate domain.Rate, date time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.latest[code]; ok {
		return domain.ErrAlreadyExists
	}

	r.latest[code] = domain.Currency{
		Code:     code,
		Rate:     rate,
		RateDate: date,
	}
	return nil
}

// update one currency
func (r *CurrencyRepoInMemory) UpdateOne(ctx context.Context, code domain.CurrencyCode, rate domain.Rate, date time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.latest[code]; !ok {
		return domain.ErrNotFound
	}

	if _, ok := r.latest[code]; !ok {
		return domain.ErrNotFound
	}

	r.latest[code] = domain.Currency{
		Code:     code,
		Rate:     rate,
		RateDate: date,
	}
	return nil
}

// Update
func (r *CurrencyRepoInMemory) UpdateAll(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for code, currency := range r.latest {
		currency.RateDate = now
		r.latest[code] = currency
	}

	return nil
}

// Delete
func (r *CurrencyRepoInMemory) DeleteAll(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.latest = make(map[domain.CurrencyCode]domain.Currency)
	return nil
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
	return r.GetOne(ctx, code)
}

func (r *CurrencyRepoInMemory) GetAtDate(
	ctx context.Context,
	code domain.CurrencyCode,
	date time.Time) (domain.Currency, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	h := r.history[code]
	for _, c := range h {
		if c.RateDate.Equal(date) {
			return c, nil
		}
	}
	return domain.Currency{}, domain.ErrNotFound
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
