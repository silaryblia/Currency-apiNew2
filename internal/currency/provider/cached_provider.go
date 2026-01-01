package provider

import (
	"Currency-apiNew2/internal/currency/domain"
	"context"
	"sync"
	"time"
)

type CachedProvider struct {
	provider  domain.RatesProvider
	cache     map[string]float64
	cacheDate time.Time
	expires   time.Time
	ttl       time.Duration
	mu        sync.RWMutex
}

func NewCachedProvider(p domain.RatesProvider, ttl time.Duration) *CachedProvider {
	return &CachedProvider{
		provider: p,
		ttl:      ttl,
	}
}

func (c *CachedProvider) GetRates(ctx context.Context) (map[string]float64, time.Time, error) {
	c.mu.RLock()
	if time.Now().Before(c.expires) && c.cache != nil {
		defer c.mu.RUnlock()
		return c.cache, c.cacheDate, nil
	}
	c.mu.RUnlock()

	return c.ForceRefresh(ctx)
}

func (c *CachedProvider) ForceRefresh(ctx context.Context) (map[string]float64, time.Time, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	rates, date, err := c.provider.GetRates(ctx)
	if err != nil {
		return nil, time.Time{}, err
	}

	c.cache = rates
	c.cacheDate = date
	c.expires = time.Now().Add(c.ttl)

	return rates, date, nil
}
