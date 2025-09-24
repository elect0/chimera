package api

import (
	"log/slog"
	"net/http"

	"github.com/elect0/chimera/internal/ports"
)

type Handler struct {
	service ports.TransformationService
	log *slog.Logger
}

func NewHandler(service ports.TransformationService, log *slog.Logger) *Handler {
	return &Handler{
		service: service,
		log: log,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", h.handleHealthCheck)
	mux.HandleFunc("/transform", h.handleImageTransformation)
}
