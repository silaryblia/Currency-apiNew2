package domain

import (
	"context"
	"time"
)

type RatesProvider interface {
	GetRates(ctx context.Context) (map[string]float64, time.Time, error)
	ForceRefresh(ctx context.Context) (map[string]float64, time.Time, error)
}
