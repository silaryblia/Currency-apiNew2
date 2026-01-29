package grpc

import (
	"Currency-apiNew2/internal/currency/service"
	"context"

	"google.golang.org/grpc/health/grpc_health_v1"
)

type HealthServer struct {
	grpc_health_v1.UnimplementedHealthServer
	svc *service.CurrencyService
}

func NewHealthServer(svc *service.CurrencyService) *HealthServer {
	return &HealthServer{svc: svc}
}

func (h *HealthServer) Check(
	ctx context.Context,
	req *grpc_health_v1.HealthCheckRequest,
) (*grpc_health_v1.HealthCheckResponse, error) {
	if h.svc.IsReady() {
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_SERVING,
		}, nil
	}

	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
	}, nil
}

func (h *HealthServer) Watch(
	req *grpc_health_v1.HealthCheckRequest,
	stream grpc_health_v1.Health_WatchServer,
) error {
	return nil
}
