package gateway

import (
	"Currency-apiNew2/internal/currency/domain"
	"Currency-apiNew2/internal/currency/service"
	"Currency-apiNew2/internal/metrics"
	"context"
	"time"

	"go.uber.org/zap"
)

type currencyGateway struct {
	service *service.CurrencyService
	logger  *zap.Logger
}

func NewCurrencyGateway(
	svc *service.CurrencyService,
	logger *zap.Logger,
) CurrencyGateway {
	return &currencyGateway{service: svc, logger: logger}
}

func (g *currencyGateway) GetLatest(
	ctx context.Context,
	code string) (domain.Currency, error) {

	c, err := domain.ParseCurrencyCode(code)
	if err != nil {
		return domain.Currency{}, err
	}
	return g.service.GetLatest(ctx, c)
}

func (g *currencyGateway) GetAtDate(
	ctx context.Context,
	code string,
	date time.Time) (domain.Currency, error) {

	c, err := domain.ParseCurrencyCode(code)
	if err != nil {
		return domain.Currency{}, err
	}
	return g.service.GetAtDate(ctx, c, date)
}

func (g *currencyGateway) GetRange(
	ctx context.Context,
	code string,
	from, to time.Time) ([]domain.Currency, error) {

	if from.After(to) {
		return nil, domain.ErrInvalidArg
	}
	c, err := domain.ParseCurrencyCode(code)
	if err != nil {
		return nil, err
	}

	return g.service.GetRange(ctx, c, from, to)
}

func (g *currencyGateway) GetAll(
	ctx context.Context) (map[domain.CurrencyCode]domain.Currency, error) {
	if !g.service.IsReady() {
		return nil, domain.ErrServiceNotReady
	}
	return g.service.GetAll(ctx)
}

func (g *currencyGateway) SyncRates(ctx context.Context) (err error) {
	defer func() {
		if r := recover(); r != nil {
			g.logger.Error("panic in SyncRates", zap.Any("panic", r))
			metrics.Ready.Set(0)
			err = domain.ErrInternal
		}
	}()

	if err = g.service.SyncRatesWithRetry(ctx); err != nil {
		metrics.Ready.Set(0)
		return err
	}

	g.service.SetReady()
	metrics.Ready.Set(1)

	return nil
}

func (g *currencyGateway) IsReady() bool {
	return g.service.IsReady()
}
