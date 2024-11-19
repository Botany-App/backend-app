package middleware

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	TotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_total_requests",
			Help: "Total de requisições recebidas pela API",
		},
		[]string{"endpoint"},
	)

	BlockedRequestsRateLimit = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_blocked_requests_rate_limit",
			Help: "Requisições bloqueadas por exceder o limite de taxa",
		},
		[]string{"endpoint"},
	)

	BlockedRequestsJail = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_blocked_requests_jail",
			Help: "Requisições bloqueadas por estar na Jail",
		},
		[]string{"endpoint"},
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "api_request_duration_seconds",
			Help:    "Duração das requisições em segundos",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint"},
	)
)

func InitPrometheus() {
	prometheus.MustRegister(TotalRequests, BlockedRequestsRateLimit, BlockedRequestsJail, RequestDuration)
}

func PrometheusHandler() http.Handler {
	return promhttp.Handler()
}
