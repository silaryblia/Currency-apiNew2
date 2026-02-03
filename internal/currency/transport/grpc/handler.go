package grpc

import (
	"Currency-apiNew2/internal/currency/domain"
	pb "Currency-apiNew2/internal/currency/proto"
	"Currency-apiNew2/internal/currency/service"
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CurrencyHandler struct {
	pb.UnimplementedCurrencyServiceServer
	service *service.CurrencyService
	logger  *zap.Logger
}

func NewCurrencyHandler(
	svc *service.CurrencyService,
	logger *zap.Logger,
) *CurrencyHandler {
	return &CurrencyHandler{
		service: svc,
		logger:  logger,
	}
}

func mapError(err error) error {
	switch err {
	case domain.ErrNotFound:
		return status.Error(codes.NotFound, err.Error())
	case domain.ErrAlreadyExists:
		return status.Error(codes.AlreadyExists, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}

func (h *CurrencyHandler) GetAtDate(
	ctx context.Context,
	req *pb.GetAtDateRequest,
) (*pb.GetOneResponse, error) {

	code, err := domain.ParseCurrencyCode(req.Code)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	c, err := h.service.GetAtDate(ctx, code, req.Date.AsTime())
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

	code, err := domain.ParseCurrencyCode(req.Code)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	items, err := h.service.GetRange(
		ctx,
		code,
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
