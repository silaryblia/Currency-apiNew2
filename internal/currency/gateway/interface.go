package gateway

import (
	"Currency-apiNew2/internal/currency/domain"
	"context"
	"time"
)

type CurrencyGateway interface {
	GetLatest(ctx context.Context, code string) (domain.Currency, error)
	GetAtDate(ctx context.Context, code string, date time.Time) (domain.Currency, error)
	GetRange(ctx context.Context, code string, from, to time.Time) ([]domain.Currency, error)
	GetAll(ctx context.Context) (map[domain.CurrencyCode]domain.Currency, error)
	SyncRates(ctx context.Context) error
	IsReady() bool
}
