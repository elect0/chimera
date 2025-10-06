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
	path := query.Get("path")
	remoteURL := query.Get("url")

	var imagePath string

	if path != "" {
		imagePath = path
	} else if remoteURL != "" {
		imagePath = remoteURL
	} else {
		http.Error(w, "invalid image path parameter", http.StatusBadRequest)
		return
	}

	width, _ := strconv.Atoi(query.Get("width"))
	height, _ := strconv.Atoi(query.Get("height"))
	quality, _ := strconv.Atoi(query.Get("quality"))

	if width <= 0 && height <= 0 {
		http.Error(w, "at least one of 'width' or 'height' parameters is invalid", http.StatusBadRequest)
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
