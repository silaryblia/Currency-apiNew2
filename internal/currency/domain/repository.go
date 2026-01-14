package domain

import (
	"context"
	"time"
)

type CurrencyRepository interface {
	GetOne(ctx context.Context, code CurrencyCode) (Currency, error)
	GetAll(ctx context.Context) (map[CurrencyCode]Currency, error)
	Create(ctx context.Context, code CurrencyCode, rate Rate, date time.Time) error
	UpdateOne(ctx context.Context, code CurrencyCode, rate Rate, date time.Time) error
	UpdateAll(ctx context.Context) error
	DeleteAll(ctx context.Context) error

	Upsert(ctx context.Context, code CurrencyCode, rate Rate, date time.Time) error
}
