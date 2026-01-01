package domain

import (
	"context"
	"time"
)

type CurrencyRepository interface {
	//GetOne(code string) (float64, error)
	//GetAll() (map[string]float64, error)

	//GetAll() ([]Currency, error)

	GetOne(ctx context.Context, code string) (Currency, error)
	GetAll(ctx context.Context) (map[string]Currency, error)
	Create(ctx context.Context, code string, rate float64, date time.Time) error
	UpdateOne(ctx context.Context, code string, rate float64, date time.Time) error
	UpdateAll(ctx context.Context) error
	DeleteAll(ctx context.Context) error

	Upsert(ctx context.Context, code string, rate float64, date time.Time) error
}
