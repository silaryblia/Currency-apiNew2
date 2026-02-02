package provider

import (
	"Currency-apiNew2/internal/currency/domain"
	"context"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

type CachedProvider struct {
	provider domain.RatesProvider

	rates     map[string]float64
	ratesDate time.Time // дата курсов
	expiresAt time.Time // когда истекает кеш
	ttl       time.Duration

	mu sync.RWMutex
	sf singleflight.Group
}

func NewCachedProvider(p domain.RatesProvider, ttl time.Duration) *CachedProvider {
	return &CachedProvider{
		provider: p,
		ttl:      ttl,
	}
}

func (c *CachedProvider) GetRates(ctx context.Context) (map[string]float64, time.Time, error) {
	c.mu.RLock()
	if time.Now().Before(c.expiresAt) && c.rates != nil {
		rates := c.rates
		date := c.ratesDate
		c.mu.RUnlock()
		return rates, date, nil
	}
	c.mu.RUnlock()

	v, err, _ := c.sf.Do("rates", func() (any, error) {
		rates, date, err := c.provider.GetRates(ctx)
		if err != nil {
			return nil, err
		}

		c.mu.Lock()
		c.rates = rates
		c.ratesDate = date
		c.expiresAt = time.Now().Add(c.ttl)
		c.mu.Unlock()

		return struct {
			rates map[string]float64
			date  time.Time
		}{rates, date}, nil
	})

	if err != nil {
		return nil, time.Time{}, err
	}

	res := v.(struct {
		rates map[string]float64
		date  time.Time
	})

	return res.rates, res.date, nil
}

func (c *CachedProvider) ForceRefresh(ctx context.Context) (map[string]float64, time.Time, error) {
	v, err, _ := c.sf.Do("rates", func() (any, error) {
		rates, date, err := c.provider.GetRates(ctx)
		if err != nil {
			return nil, err
		}

		c.mu.Lock()
		c.rates = rates
		c.ratesDate = date
		c.expiresAt = time.Now().Add(c.ttl)
		c.mu.Unlock()

		return struct {
			rates map[string]float64
			date  time.Time
		}{rates, date}, nil
	})

	if err != nil {
		return nil, time.Time{}, err
	}

	res := v.(struct {
		rates map[string]float64
		date  time.Time
	})

	return res.rates, res.date, nil
}
