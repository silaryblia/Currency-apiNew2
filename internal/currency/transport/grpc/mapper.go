package grpc

import (
	"Currency-apiNew2/internal/currency/domain"
	pb "Currency-apiNew2/internal/currency/proto"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToProtoCurrency(c domain.Currency) *pb.Currency {
	return &pb.Currency{
		Code:     c.Code.String(),
		Rate:     c.Rate.Float64(),
		RateDate: timestamppb.New(c.RateDate),
	}
}

func ToProtoCurrencyMap(src map[domain.CurrencyCode]domain.Currency) map[string]*pb.Currency {
	dst := make(map[string]*pb.Currency, len(src))
	for code, c := range src {
		dst[code.String()] = ToProtoCurrency(c)
	}
	return dst
}
