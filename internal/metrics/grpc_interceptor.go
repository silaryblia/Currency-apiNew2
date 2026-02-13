package metrics

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

		start := time.Now()
		resp, err := handler(ctx, req)
		elapsed := time.Since(start)

		st := status.Code(err).String()

		GrpcRequests.WithLabelValues(info.FullMethod, st).Inc()
		GrpcLatency.WithLabelValues(info.FullMethod).Observe(elapsed.Seconds())

		return resp, err
	}
}
