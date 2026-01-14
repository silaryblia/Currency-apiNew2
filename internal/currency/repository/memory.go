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
	mu     sync.RWMutex
	data   map[domain.CurrencyCode]domain.Currency
	logger *zap.Logger
}

type currencyRecord struct {
	rate float64
	date time.Time
}

func NewCurrencyRepoInMemory(logger *zap.Logger) *CurrencyRepoInMemory {
	return &CurrencyRepoInMemory{
		data:   make(map[domain.CurrencyCode]domain.Currency),
		logger: logger,
	}
}

func (r *CurrencyRepoInMemory) GetOne(ctx context.Context, code domain.CurrencyCode) (domain.Currency, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	c, ok := r.data[code]
	if !ok {
		return domain.Currency{}, domain.ErrNotFound
	}

	return c, nil
}

func (r *CurrencyRepoInMemory) GetAll(ctx context.Context) (map[domain.CurrencyCode]domain.Currency, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	res := make(map[domain.CurrencyCode]domain.Currency, len(r.data))
	for k, v := range r.data {
		res[k] = v
	}

	return res, nil
}

func (r *CurrencyRepoInMemory) Upsert(ctx context.Context, code domain.CurrencyCode, rate domain.Rate, date time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[code] = domain.Currency{
		Code:     code,
		Rate:     rate,
		RateDate: date,
	}
	return nil
}

// add new currency
func (r *CurrencyRepoInMemory) Create(ctx context.Context, code domain.CurrencyCode, rate domain.Rate, date time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[code]; ok {
		return domain.ErrAlreadyExists
	}

	r.data[code] = domain.Currency{
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

	if _, ok := r.data[code]; !ok {
		return domain.ErrNotFound
	}

	if _, ok := r.data[code]; !ok {
		return domain.ErrNotFound
	}

	r.data[code] = domain.Currency{
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
	for code, currency := range r.data {
		currency.RateDate = now
		r.data[code] = currency
	}

	return nil
}

// Delete
func (r *CurrencyRepoInMemory) DeleteAll(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data = make(map[domain.CurrencyCode]domain.Currency)
	return nil
}
