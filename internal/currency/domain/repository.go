package domain

import (
	"context"
	"time"
)

type CurrencyRepository interface {
	GetAll(ctx context.Context) (map[CurrencyCode]Currency, error)

	SaveRate(ctx context.Context, code CurrencyCode, rate Rate, date time.Time) error
	GetLatest(ctx context.Context, code CurrencyCode) (Currency, error)
	GetAtDate(ctx context.Context, code CurrencyCode, date time.Time) (Currency, error)
	GetRange(ctx context.Context, code CurrencyCode, from, to time.Time) ([]Currency, error)
}
