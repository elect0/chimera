package api

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/elect0/chimera/internal/domain"
)

func (h *Handler) handleImageTransformation(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	query := r.URL.Query()
	imagePath := query.Get("path")
	width, _ := strconv.Atoi(query.Get("width"))
	height, _ := strconv.Atoi(query.Get("height"))
	quality, _ := strconv.Atoi(query.Get("quality"))

	if imagePath == "" {
		http.Error(w, "invalid image path parameter", http.StatusBadRequest)
		return
	}

	if width <= 0 {
		http.Error(w, "invalid width parameter", http.StatusBadRequest)
		return
	}

	opts := domain.TransformationOptions{
		Width:   width,
		Height:  height,
		Quality: quality,
	}

	processedImage, err := h.service.Process(r.Context(), opts, imagePath)
	if err != nil {
		http.Error(w, "failed to process image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(processedImage)

	h.log.Info("request processed successfully", slog.Duration("duration", time.Since(start)), slog.Int("status", http.StatusOK), slog.String("path", r.URL.Path))
}
