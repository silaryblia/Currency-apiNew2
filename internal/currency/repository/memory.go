package repository

import (
	"Currency-apiNew2/internal/currency/domain"
	"context"
	_ "math/rand"
	"strings"
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
	data   map[string]domain.Currency
	logger *zap.Logger
}

type currencyRecord struct {
	rate float64
	date time.Time
}

//func NewCurrencyRepoInMemory(logger *zap.Logger) *CurrencyRepoInMemory {
//	now := time.Now()

//	return &CurrencyRepoInMemory{
//		data: map[string]domain.Currency{
//			"USD": {Code: "USD", Rate: 80, RateDate: now},
//			"EUR": {Code: "EUR", Rate: 85, RateDate: now},
//			"AED": {Code: "AED", Rate: 20, RateDate: now},
//		},
//		logger: logger,
//	}
//}

func NewCurrencyRepoInMemory(logger *zap.Logger) *CurrencyRepoInMemory {
	return &CurrencyRepoInMemory{
		data:   make(map[string]domain.Currency),
		logger: logger,
	}
}

func (r *CurrencyRepoInMemory) GetOne(ctx context.Context, code string) (domain.Currency, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	code = strings.ToUpper(strings.TrimSpace(code))

	c, ok := r.data[code]
	if !ok {
		return domain.Currency{}, domain.ErrNotFound
	}

	//r.logger.Debug("repo: Get success", zap.String("code", code), zap.Float64("rate", rec.rate))
	return c, nil
}

func (r *CurrencyRepoInMemory) GetAll(ctx context.Context) (map[string]domain.Currency, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	res := make(map[string]domain.Currency, len(r.data))
	for k, v := range r.data {
		res[k] = v
	}

	return res, nil
}

func (r *CurrencyRepoInMemory) Upsert(ctx context.Context, code string, rate float64, date time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	code = strings.ToUpper(code)
	r.data[code] = domain.Currency{
		Code:     code,
		Rate:     rate,
		RateDate: date,
	}
	return nil
}

// add new currency
func (r *CurrencyRepoInMemory) Create(ctx context.Context, code string, rate float64, date time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	code = strings.ToUpper(strings.TrimSpace(code))

	if _, ok := r.data[code]; ok {
		return domain.ErrAlreadyExists
	}

	//r.data[code] = domain.Currency{rate: rate, date: date}
	r.data[code] = domain.Currency{
		Code:     code,
		Rate:     rate,
		RateDate: date,
	}
	return nil
}

// update one currency
func (r *CurrencyRepoInMemory) UpdateOne(ctx context.Context, code string, rate float64, date time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	code = strings.ToUpper(strings.TrimSpace(code))

	if _, ok := r.data[code]; !ok {
		return domain.ErrNotFound
	}

	//r.data[code] = rate

	cur, ok := r.data[code]
	if !ok {
		return domain.ErrNotFound
	}

	cur.Rate = rate
	cur.RateDate = date
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

	// clear map
	r.data = make(map[string]domain.Currency)
	return nil
}
