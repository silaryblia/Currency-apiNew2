package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	ProviderRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "currency",
			Subsystem: "provider",
			Name:      "requests_total",
			Help:      "Total requests to external currency provider",
		},
		[]string{"provider"},
	)

	ProviderErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "currency",
			Subsystem: "provider",
			Name:      "errors_total",
			Help:      "Total errors from external currency provider",
		},
		[]string{"provider", "reason"},
	)

	ProviderLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "currency",
			Subsystem: "provider",
			Name:      "latency_seconds",
			Help:      "External provider request latency",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"provider"},
	)
)
