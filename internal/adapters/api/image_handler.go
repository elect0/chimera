package api

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/elect0/chimera/internal/domain"
	"github.com/h2non/bimg"
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
	cropStrategy := query.Get("crop")

	targetType := h.negotiateFormat(r)

	if width <= 0 && height <= 0 {
		http.Error(w, "at least one of 'width' or 'height' parameters is invalid", http.StatusBadRequest)
		return
	}

	opts := domain.TransformationOptions{
		Width:      width,
		Height:     height,
		Quality:    quality,
		Crop:       cropStrategy,
		TargetType: targetType,
	}

	processedImage, err := h.service.Process(r.Context(), opts, imagePath)
	if err != nil {
		http.Error(w, "failed to process image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/"+bimg.ImageTypeName(targetType))
	w.Write(processedImage)

	h.log.Info("request processed successfully", slog.Duration("duration", time.Since(start)), slog.Int("status", http.StatusOK), slog.String("path", r.URL.Path))
}

func (h *Handler) negotiateFormat(r *http.Request) bimg.ImageType {
	acceptHeader := r.Header.Get("Accept")

	if strings.Contains(acceptHeader, "image/avif") {
		h.log.Debug("client supports AVIF, selecting AVIF")
		return bimg.AVIF
	}

	if strings.Contains(acceptHeader, "image/webp") {
		h.log.Debug("client supports WEBP, selecting WEBP")
		return bimg.WEBP
	}

	h.log.Debug("client doesn't support neither WEBP nor AVIF, falling back to JPEG")
	return bimg.JPEG
}
