package grpc

import (
	"Currency-apiNew2/internal/currency/gateway"
	pb "Currency-apiNew2/internal/currency/proto"
	"context"

	"go.uber.org/zap"
)

type CurrencyHandler struct {
	pb.UnimplementedCurrencyServiceServer
	gateway gateway.CurrencyGateway
	logger  *zap.Logger
}

func NewCurrencyHandler(
	gw gateway.CurrencyGateway,
	logger *zap.Logger,
) *CurrencyHandler {
	return &CurrencyHandler{
		gateway: gw,
		logger:  logger,
	}
}

func (h *CurrencyHandler) GetAtDate(
	ctx context.Context,
	req *pb.GetAtDateRequest,
) (*pb.GetOneResponse, error) {

	c, err := h.gateway.GetAtDate(ctx, req.Code, req.Date.AsTime())
	if err != nil {
		return nil, mapError(err)
	}

	return &pb.GetOneResponse{
		Currency: ToProtoCurrency(c),
	}, nil
}

func (h *CurrencyHandler) GetRange(
	ctx context.Context,
	req *pb.GetRangeRequest,
) (*pb.GetRangeResponse, error) {

	items, err := h.gateway.GetRange(
		ctx,
		req.Code,
		req.From.AsTime(),
		req.To.AsTime(),
	)
	if err != nil {
		return nil, mapError(err)
	}

	res := make([]*pb.Currency, 0, len(items))
	for _, c := range items {
		res = append(res, ToProtoCurrency(c))
	}

	return &pb.GetRangeResponse{Currencies: res}, nil
}

func (h *CurrencyHandler) GetLatest(
	ctx context.Context,
	req *pb.GetOneRequest,
) (*pb.GetOneResponse, error) {

	c, err := h.gateway.GetLatest(ctx, req.Code)
	if err != nil {
		return nil, mapError(err)
	}

	return &pb.GetOneResponse{
		Currency: ToProtoCurrency(c),
	}, nil
}

func (h *CurrencyHandler) GetAll(
	ctx context.Context,
	_ *pb.GetAllRequest,
) (*pb.GetAllResponse, error) {

	items, err := h.gateway.GetAll(ctx)
	if err != nil {
		return nil, mapError(err)
	}

	return &pb.GetAllResponse{
		Currencies: ToProtoCurrencyMap(items),
	}, nil
}

func (h *CurrencyHandler) SyncRates(
	ctx context.Context,
	_ *pb.GetAllRequest,
) (*pb.GetAllResponse, error) {

	if err := h.gateway.SyncRates(ctx); err != nil {
		return nil, mapError(err)
	}

	items, err := h.gateway.GetAll(ctx)
	if err != nil {
		return nil, mapError(err)
	}

	return &pb.GetAllResponse{
		Currencies: ToProtoCurrencyMap(items),
	}, nil
}
