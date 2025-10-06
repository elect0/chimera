package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/elect0/chimera/internal/config"
	"github.com/elect0/chimera/internal/ports"
)

type HTTPOriginRepository struct {
	client *http.Client
	cfg    *config.Config
	log    *slog.Logger
}

func NewHTTPOriginRepository(cfg *config.Config, log *slog.Logger) *HTTPOriginRepository {
	client := &http.Client{
		Timeout: 30,
	}

	return &HTTPOriginRepository{
		client: client,
		cfg:    cfg,
		log:    log,
	}
}

func (r *HTTPOriginRepository) Get(ctx context.Context, imageURL string) ([]byte, error) {
	log := r.log.With(slog.String("imageURL", imageURL))
	log.Debug("fetching image from remote url")

	if err := r.isPubliclyRoutable(imageURL); err != nil {
		log.Warn("ssrf attempt rejected", slog.String("error", err.Error()))
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, imageURL, nil)
	if err != nil {
		log.Error("failed to create http request", slog.String("error", err.Error()))
		return nil, err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		log.Error("failed to fetch remote url", slog.String("error", err.Error()))
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("remote server returned status code %d", resp.StatusCode)
	}

	maxSizeBytes := int64(r.cfg.Security.RemoteFetch.MaxDownloadSizeMB) * 1024
	if resp.ContentLength > maxSizeBytes {
		return nil, fmt.Errorf("remote file size (%d bytes) exceeds max limit", resp.ContentLength)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return nil, fmt.Errorf("invalid content type '%s' must be an image", contentType)
	}

	limitedReader := &io.LimitedReader{R: resp.Body, N: maxSizeBytes}
	body, err := io.ReadAll(limitedReader)
	if err != nil {
		log.Error("failed to read response body", slog.String("error", err.Error()))
		return nil, err
	}

	log.Debug("successfully fetched image from remote url", slog.Int("sizes_bytes", len(body)))
	return body, nil
}

func (r *HTTPOriginRepository) isPubliclyRoutable(imageURL string) error {
	parsedURL, err := url.Parse(imageURL)
	if err != nil {
		return fmt.Errorf("invalid url: %w", err)
	}

	ips, err := net.LookupIP(parsedURL.Hostname())
	if err != nil {
		return fmt.Errorf("dns lookup failed: %w", err)
	}

	for _, ip := range ips {
		if ip.IsPrivate() || ip.IsLoopback() || ip.IsUnspecified() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() {
			return errors.New("url resolves to a non-public ip address")
		}
	}

	return nil
}

var _ ports.OriginRepository = (*HTTPOriginRepository)(nil)
