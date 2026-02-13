package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	GrpcRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "currency",
			Subsystem: "grpc",
			Name:      "requests_total",
			Help:      "Total gRPC requests",
		},
		[]string{"method", "status"},
	)

	GrpcLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "currency",
			Subsystem: "grpc",
			Name:      "latency_seconds",
			Help:      "gRPC latency",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method"},
	)
)

func MustRegister() {
	prometheus.MustRegister(
		GrpcRequests,
		GrpcLatency,
		Ready,

		ProviderRequests,
		ProviderErrors,
		ProviderLatency,
	)
}
