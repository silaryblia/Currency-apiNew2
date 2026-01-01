package grpc

import (
	_ "Currency-apiNew2/internal/currency/domain"
	pb "Currency-apiNew2/internal/currency/proto"
	"Currency-apiNew2/internal/currency/service"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
)

type CurrencyServer struct {
	pb.UnimplementedCurrencyServiceServer
	service *service.CurrencyService
}

func NewCurrencyServer(service *service.CurrencyService) *CurrencyServer {
	return &CurrencyServer{service: service}
}

func (s *CurrencyServer) GetAll(ctx context.Context, req *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	currencies, err := s.service.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var currenciesMap = make(map[string]*pb.Currency)
	for code, c := range currencies {
		currenciesMap[code] = &pb.Currency{
			Code:     c.Code,
			Rate:     c.Rate,
			RateDate: timestamppb.New(c.RateDate),
		}
	}

	return &pb.GetAllResponse{Currencies: currenciesMap}, nil
}

func (s *CurrencyServer) GetOne(ctx context.Context, req *pb.GetOneRequest) (*pb.GetOneResponse, error) {
	currency, err := s.service.GetOne(ctx, req.Code)
	if err != nil {
		return nil, err
	}

	return &pb.GetOneResponse{
		Currency: &pb.Currency{
			Code:     currency.Code,
			Rate:     currency.Rate,
			RateDate: timestamppb.New(currency.RateDate),
		},
	}, nil
}

func (s *CurrencyServer) Create(ctx context.Context, req *pb.Currency) (*pb.Currency, error) {
	err := s.service.Create(ctx, req.Code, req.Rate, req.RateDate.AsTime())
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (s *CurrencyServer) UpdateOne(ctx context.Context, req *pb.Currency) (*pb.Currency, error) {
	err := s.service.UpdateOne(ctx, req.Code, req.Rate, req.RateDate.AsTime())
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (s *CurrencyServer) SyncRates(ctx context.Context, req *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	err := s.service.SyncRates(ctx)
	if err != nil {
		return nil, err
	}

	return s.GetAll(ctx, req)
}

func RunServer(service *service.CurrencyService) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCurrencyServiceServer(grpcServer, NewCurrencyServer(service))

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
