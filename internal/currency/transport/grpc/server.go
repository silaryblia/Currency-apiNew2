package grpc

import (
	"Currency-apiNew2/internal/currency/domain"
	pb "Currency-apiNew2/internal/currency/proto"
	"Currency-apiNew2/internal/currency/service"
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type CurrencyServer struct {
	pb.UnimplementedCurrencyServiceServer
	service *service.CurrencyService
}

func NewCurrencyServer(service *service.CurrencyService) *CurrencyServer {
	return &CurrencyServer{service: service}
}

func (s *CurrencyServer) GetAll(
	ctx context.Context,
	req *pb.GetAllRequest,
) (*pb.GetAllResponse, error) {

	currencies, err := s.service.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.GetAllResponse{Currencies: ToProtoCurrencyMap(currencies)}, nil
}

func (s *CurrencyServer) SyncRates(ctx context.Context, req *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	err := s.service.SyncRates(ctx)
	if err != nil {
		return nil, err
	}

	return s.GetAll(ctx, req)
}

func (s *CurrencyServer) GetLatest(
	ctx context.Context,
	req *pb.GetOneRequest,
) (*pb.GetOneResponse, error) {

	code, err := domain.ParseCurrencyCode(req.Code)
	if err != nil {
		return nil, err
	}

	currency, err := s.service.GetLatest(ctx, code)
	if err != nil {
		return nil, err
	}

	return &pb.GetOneResponse{
		Currency: ToProtoCurrency(currency),
	}, nil
}

func RunServer(
	ctx context.Context,
	service *service.CurrencyService,
	port string) error {
	addr := ":" + port

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("grpc listen %s: %w", addr, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCurrencyServiceServer(grpcServer, NewCurrencyServer(service))
	grpc_health_v1.RegisterHealthServer(grpcServer, NewHealthServer(service))
	reflection.Register(grpcServer)

	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()

	fmt.Println("gRPC listening on", addr)
	return grpcServer.Serve(lis)
}

func (s *CurrencyServer) GetAtDate(
	ctx context.Context,
	req *pb.GetAtDateRequest,
) (*pb.GetOneResponse, error) {
	code, err := domain.ParseCurrencyCode(req.Code)
	if err != nil {
		return nil, err
	}

	date := req.Date.AsTime()

	currency, err := s.service.GetAtDate(ctx, code, date)
	if err != nil {
		return nil, err
	}

	return &pb.GetOneResponse{
		Currency: ToProtoCurrency(currency),
	}, nil
}
