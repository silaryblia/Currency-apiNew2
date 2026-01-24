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
	durationMs time.Duration) {
	n.logger.Warn(
		"slow operation detected",
		zap.String("operation", name),
		zap.Duration("duration_ms", durationMs))
}

func (n *LoggerNotificationService) RateSpike(
	ctx context.Context,
	code domain.CurrencyCode,
	oldRate, newRate domain.Rate) {

	changePercent := (float64(newRate-oldRate) / float64(oldRate)) * 100

	n.logger.Warn(
		"rate spike detected",
		zap.String("code", code.String()),
		zap.Float64("old_rate", float64(oldRate)),
		zap.Float64("new_rate", float64(newRate)),
		zap.Float64("change_percent", changePercent))
}
