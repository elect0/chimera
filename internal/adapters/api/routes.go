package api

import (
	"log/slog"
	"net/http"

	"github.com/elect0/chimera/internal/config"
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

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", h.handleHealthCheck)

	transformHandler := http.HandlerFunc(h.handleImageTransformation)

	mux.Handle("/transform", h.SignatureMiddleware(transformHandler))
}
