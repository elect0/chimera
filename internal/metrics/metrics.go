package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HTTPRequestTotals = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "chimera_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"status", "method", "path"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "chimera_http_request_duration_seconds",
			Help:    "Latency of HTTP request in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"status", "method", "path"},
	)

	CacheHitTotals = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "chimera_cache_hits_total",
			Help: "Total number of cache hits",
		},
	)

	CacheMissesTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "chimera_cache_misses_total",
			Help: "Total number of cache misses",
		},
	)
)
