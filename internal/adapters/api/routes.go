package api

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/elect0/chimera/internal/config"
	"github.com/elect0/chimera/internal/metrics"
	"github.com/elect0/chimera/internal/ports"
)

type Handler struct {
	service ports.TransformationService
	log     *slog.Logger
	cfg     *config.Config
}

func NewHandler(service ports.TransformationService, log *slog.Logger, cfg *config.Config) *Handler {
	return &Handler{
		service: service,
		log:     log,
		cfg:     cfg,
	}
}

func (h *Handler) MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		d := &responseData{
			ResponseWriter: w,
		}

		next.ServeHTTP(d, r)

		duration := time.Since(start)

		metrics.HTTPRequestDuration.WithLabelValues(
			strconv.Itoa(d.status),
			r.Method,
			r.URL.Path,
		).Observe(duration.Seconds())

		metrics.HTTPRequestTotals.WithLabelValues(
			strconv.Itoa(d.status),
			r.Method,
			r.URL.Path,
		).Inc()
	})
}

type responseData struct {
	http.ResponseWriter
	status int
}

func (r *responseData) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	healthHandler := http.HandlerFunc(h.handleHealthCheck)
	transformHandler := http.HandlerFunc(h.handleImageTransformation)

	mux.Handle("/health", h.MetricsMiddleware(healthHandler))

	if h.cfg.Security.HMACEnabled {
		h.log.Info("HMAC Signature validation is enabled for /transform")
		mux.Handle("/transform", h.MetricsMiddleware(h.SignatureMiddleware(transformHandler)))
	} else {
		h.log.Info("HMAC Signature validation is disabled for /transform")
		mux.Handle("/transform", h.MetricsMiddleware(transformHandler))
	}

}
