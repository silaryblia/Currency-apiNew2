package grpc

import (
	"Currency-apiNew2/internal/currency/gateway"
	"context"

	//"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	//"google.golang.org/grpc/status"
)

type HealthServer struct {
	//grpc_health_v1.UnimplementedHealthServer
	healthpb.UnimplementedHealthServer
	gateway gateway.CurrencyGateway
}

func NewHealthServer(gw gateway.CurrencyGateway) *HealthServer {
	return &HealthServer{gateway: gw}
}

func (s *HealthServer) Check(
	ctx context.Context,
	req *healthpb.HealthCheckRequest,
) (*healthpb.HealthCheckResponse, error) {

	if s.gateway.IsReady() {
		return &healthpb.HealthCheckResponse{
			Status: healthpb.HealthCheckResponse_SERVING,
		}, nil
	}

	return &healthpb.HealthCheckResponse{
		Status: healthpb.HealthCheckResponse_NOT_SERVING,
	}, nil
}
