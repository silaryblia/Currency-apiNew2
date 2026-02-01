package notification

import (
	"Currency-apiNew2/internal/currency/domain"
	"context"
	"time"

	"go.uber.org/zap"
)

type LoggerNotificationService struct {
	logger *zap.Logger
}

func NewLoggerNotificationService(logger *zap.Logger) *LoggerNotificationService {
	return &LoggerNotificationService{logger: logger}
}

func (n *LoggerNotificationService) SlowOperation(
	ctx context.Context,
	name string,
	duration time.Duration) {
	n.logger.Warn(
		"slow operation detected",
		zap.String("operation", name),
		zap.Duration("duration_ms", duration))
}

func (n *LoggerNotificationService) RateSpike(
	ctx context.Context,
	code domain.CurrencyCode,
	oldRate, newRate domain.Rate) {

	oldRateFloat := oldRate.Float64()
	newRateFloat := newRate.Float64()

	var changePercent float64
	if oldRateFloat != 0 {
		changePercent = ((newRateFloat - oldRateFloat) / oldRateFloat) * 100
	}

	n.logger.Warn(
		"rate spike detected",
		zap.String("code", string(code)), // код валюты это строка
		zap.Float64("old_rate", oldRateFloat),
		zap.Float64("new_rate", newRateFloat),
		zap.Float64("change_percent", changePercent))
}
