package transformation

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/elect0/chimera/internal/domain"
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
		return cachedImage, nil
	}

	if err != redis.Nil {
		log.Error("error getting from cache", slog.String("error", err.Error()))
	}
	log.Info("cache miss")

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

	go func() {
		err := s.cacheRepo.Set(context.Background(), cacheKey, newImage)
		if err != nil {
			log.Error("failed to set item in cache", slog.String("error", err.Error()))
		}
		log.Info("successfully set item in cache")
	}()

	return newImage, nil
}

var _ ports.TransformationService = (*Service)(nil)
