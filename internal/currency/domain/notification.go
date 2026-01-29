package domain

import (
	"context"
	"time"
)

type NotificationService interface {
	SlowOperation(ctx context.Context, name string, duration time.Duration)
	RateSpike(ctx context.Context, code CurrencyCode, oldRate, newRate Rate)
}
