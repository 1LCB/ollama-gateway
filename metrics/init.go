package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gateway_requests_total",
			Help: "Total number of HTTP requests made to the gateway",
		},
		[]string{"method", "endpoint", "status"},
	)

	RequestsDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "gateway_requests_duration_seconds",
			Help:    "Histogram of latencies for HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	RequestsByIP = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gateway_requests_by_ip_total",
			Help: "Total number of HTTP requests grouped by IP address",
		},
		[]string{"ip", "endpoint", "status"},
	)

	RequestsSuccess = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gateway_requests_success_total",
			Help: "Total number of successful HTTP requests",
		},
		[]string{"method", "endpoint"},
	)

	RequestsFailure = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gateway_requests_failure_total",
			Help: "Total number of failed HTTP requests",
		},
		[]string{"method", "endpoint"},
	)

	APIKeysTotal = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "gateway_api_keys_total",
			Help: "Total number of API keys configured in the gateway",
		},
	)

	OllamaServersTotal = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "gateway_ollama_servers_total",
			Help: "Total number of Ollama servers configured in the gateway",
		},
	)
)

func InitMetrics() {
	prometheus.MustRegister(RequestsTotal)
	prometheus.MustRegister(RequestsDuration)
	prometheus.MustRegister(RequestsByIP)
	prometheus.MustRegister(RequestsSuccess)
	prometheus.MustRegister(RequestsFailure)
	prometheus.MustRegister(APIKeysTotal)
	prometheus.MustRegister(OllamaServersTotal)
}
