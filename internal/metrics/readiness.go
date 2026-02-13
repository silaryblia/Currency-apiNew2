package metrics

import "github.com/prometheus/client_golang/prometheus"

var Ready = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Namespace: "currency",
		Name:      "ready",
		Help:      "Service readiness (1 = ready, 0 = not ready)",
	},
)
