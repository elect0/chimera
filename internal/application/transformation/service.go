package transformation

import (
	"context"
	"log/slog"

	"github.com/elect0/chimera/internal/domain"
	"github.com/elect0/chimera/internal/ports"
	"github.com/h2non/bimg"
)

type Service struct {
	log *slog.Logger
}

func NewService(log *slog.Logger) *Service {
	return &Service{
		log: log,
	}
}

func (s *Service) Process(ctx context.Context, opts domain.TransformationOptions, imageBuffer []byte) ([]byte, error) {
	s.log.Debug("processing image with options", slog.Int("width", opts.Width), slog.Int("height", opts.Height))
	
	image := bimg.NewImage(imageBuffer)

	bimgOptions := bimg.Options{
		Width: opts.Width,
		Height: opts.Height,
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
