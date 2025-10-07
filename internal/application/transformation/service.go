package transformation

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/elect0/chimera/internal/domain"
	"github.com/elect0/chimera/internal/metrics"
	"github.com/elect0/chimera/internal/ports"
	"github.com/h2non/bimg"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	log            *slog.Logger
	s3OriginRepo   ports.OriginRepository
	cacheRepo      ports.CacheRepository
	httpOriginRepo ports.OriginRepository
}

func NewService(log *slog.Logger, originRepo ports.OriginRepository, cacheRepo ports.CacheRepository, httpRepo ports.OriginRepository) *Service {
	return &Service{
		log:            log,
		s3OriginRepo:   originRepo,
		httpOriginRepo: httpRepo,
		cacheRepo:      cacheRepo,
	}
}

func (s *Service) Process(ctx context.Context, opts domain.TransformationOptions, imagePath string) ([]byte, error) {

	cacheKey := fmt.Sprintf("%s:w%d:h%d:q%d", imagePath, opts.Width, opts.Height, opts.Quality)
	log := s.log.With(slog.String("cacheKey", cacheKey))

	cachedImage, err := s.cacheRepo.Get(ctx, cacheKey)
	if err == nil {
		log.Info("cache hit")
		metrics.CacheHitTotals.Inc()
		return cachedImage, nil
	}

	if err != redis.Nil {
		log.Error("error getting from cache", slog.String("error", err.Error()))
	}
	log.Info("cache miss")
	metrics.CacheMissesTotal.Inc()

	var originalImage []byte
	if strings.HasPrefix(imagePath, "http") {
		originalImage, err = s.httpOriginRepo.Get(ctx, imagePath)
	} else {
		originalImage, err = s.s3OriginRepo.Get(ctx, imagePath)
	}

	if err != nil {
		log.Error("failed to get original image from origin", slog.String("error", err.Error()))
		return nil, err
	}

	image := bimg.NewImage(originalImage)
	bimgOptions := bimg.Options{
		Width:   opts.Width,
		Height:  opts.Height,
		Quality: opts.Quality,
		Crop:    true,
		Type:    opts.TargetType,
	}

	if opts.Crop == "smart" {
		bimgOptions.Gravity = bimg.GravitySmart
	}

	newImage, err := image.Process(bimgOptions)
	if err != nil {
		s.log.Error("failed to process image", slog.String("error", err.Error()))
		return nil, err
	}

	if opts.Watermark.Path != "" {
		s.log.Debug("watermark requested, fetching watermark image", slog.String("path", opts.Watermark.Path))

		watermarkBuffer, err := s.s3OriginRepo.Get(ctx, opts.Watermark.Path)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch watermark image: %w", err)
		}

		processedSize, err := bimg.Size(newImage)
		if err != nil {
			return nil, err
		}
		watermarkSize, err := bimg.Size(watermarkBuffer)
		if err != nil {
			return nil, err
		}

		top, left := calculateCoordinates(processedSize, watermarkSize, opts.Watermark.Position)

		imageToWatermark := bimg.NewImage(newImage)

		watermark := bimg.WatermarkImage{
			Buf:     watermarkBuffer,
			Opacity: opts.Watermark.Opacity,
			Top:     top,
			Left:    left,
		}

		newImage, err = imageToWatermark.WatermarkImage(watermark)
		if err != nil {
			s.log.Error("failed to apply watermark", slog.String("error", err.Error()))
			return nil, err
		}
	}

	go func() {
		err := s.cacheRepo.Set(context.Background(), cacheKey, newImage)
		if err != nil {
			log.Error("failed to set item in cache", slog.String("error", err.Error()))
		}
		log.Info("successfully set item in cache")
	}()

	return newImage, nil
}

func calculateCoordinates(baseSize, watermarkSize bimg.ImageSize, gravity bimg.Gravity) (top, left int) {
	switch gravity {
	case bimg.GravityNorth:
		left = (baseSize.Width - watermarkSize.Width) / 2
		top = 0
	case bimg.GravitySouth:
		left = (baseSize.Width - watermarkSize.Width) / 2
		top = baseSize.Height - watermarkSize.Height
	case bimg.GravityEast:
		left = baseSize.Width - watermarkSize.Width
		top = (baseSize.Height - watermarkSize.Height) / 2
	case bimg.GravityWest:
		left = 0
		top = (baseSize.Height - watermarkSize.Height) / 2
	case bimg.GravityCentre:
		fallthrough
	default:
		// Default to center
		left = (baseSize.Width - watermarkSize.Width) / 2
		top = (baseSize.Height - watermarkSize.Height) / 2
	}
	return top, left
}

var _ ports.TransformationService = (*Service)(nil)
