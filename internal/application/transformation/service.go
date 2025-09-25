package transformation

import (
	"context"
	"log/slog"

	"github.com/elect0/chimera/internal/domain"
	"github.com/elect0/chimera/internal/ports"
	"github.com/h2non/bimg"
)

type Service struct {
	log        *slog.Logger
	originRepo ports.OriginRepository
}

func NewService(log *slog.Logger, originRepo ports.OriginRepository) *Service {
	return &Service{
		log:        log,
		originRepo: originRepo,
	}
}

func (s *Service) Process(ctx context.Context, opts domain.TransformationOptions, imagePath string) ([]byte, error) {
	log := s.log.With(slog.String("imagePath", imagePath))
	log.Debug("processing image request")

	originalImage, err := s.originRepo.Get(ctx, imagePath)
	if err != nil {
		log.Error("failed to get original image from origin", slog.String("error", err.Error()))
		return nil, err
	}

	image := bimg.NewImage(originalImage)
	bimgOptions := bimg.Options{
		Width:   opts.Width,
		Height:  opts.Height,
		Quality: opts.Quality,
	}

	newImage, err := image.Process(bimgOptions)
	if err != nil {
		s.log.Error("failed to process image", slog.String("error", err.Error()))
		return nil, err
	}

	return newImage, nil
}

var _ ports.TransformationService = (*Service)(nil)
