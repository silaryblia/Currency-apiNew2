package grpc

import (
	"Currency-apiNew2/internal/currency/gateway"
	"context"

	"google.golang.org/grpc/health/grpc_health_v1"
)

type HealthServer struct {
	grpc_health_v1.UnimplementedHealthServer
	gateway gateway.CurrencyGateway
}

func NewHealthServer(gw gateway.CurrencyGateway) *HealthServer {
	return &HealthServer{gateway: gw}
}

func (h *HealthServer) Check(
	ctx context.Context,
	_ *grpc_health_v1.HealthCheckRequest,
) (*grpc_health_v1.HealthCheckResponse, error) {

	if !h.gateway.IsReady() {
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
		}, nil
	}

	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (h *HealthServer) Watch(
	req *grpc_health_v1.HealthCheckRequest,
	stream grpc_health_v1.Health_WatchServer,
) error {
	return nil
}
