package grpc

import (
	"Currency-apiNew2/internal/currency/gateway"
	pb "Currency-apiNew2/internal/currency/proto"
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type CurrencyServer struct {
	pb.UnimplementedCurrencyServiceServer
	gateway gateway.CurrencyGateway
}

func NewCurrencyServer(gw gateway.CurrencyGateway) *CurrencyServer {
	return &CurrencyServer{gateway: gw}
}

func (s *CurrencyServer) GetAll(
	ctx context.Context,
	_ *pb.GetAllRequest,
) (*pb.GetAllResponse, error) {

	currencies, err := s.gateway.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.GetAllResponse{Currencies: ToProtoCurrencyMap(currencies)}, nil
}

func (s *CurrencyServer) SyncRates(ctx context.Context, req *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	err := s.gateway.SyncRates(ctx)
	if err != nil {
		return nil, err
	}

	return s.GetAll(ctx, req)
}

func (s *CurrencyServer) GetLatest(
	ctx context.Context,
	req *pb.GetOneRequest,
) (*pb.GetOneResponse, error) {

	currency, err := s.gateway.GetLatest(ctx, req.Code)
	if err != nil {
		return nil, err
	}

	return &pb.GetOneResponse{
		Currency: ToProtoCurrency(currency),
	}, nil
}

func RunServer(
	ctx context.Context,
	gateway gateway.CurrencyGateway,
	port string) error {

	addr := ":" + port

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("grpc listen %s: %w", addr, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCurrencyServiceServer(grpcServer, NewCurrencyServer(gateway))
	grpc_health_v1.RegisterHealthServer(grpcServer, NewHealthServer(gateway))
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

	date := req.Date.AsTime()

	currency, err := s.gateway.GetAtDate(ctx, req.Code, date)
	if err != nil {
		return nil, err
	}

	return &pb.GetOneResponse{
		Currency: ToProtoCurrency(currency),
	}, nil
}
