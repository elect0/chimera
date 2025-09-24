package transformation

import (
	"context"
	"log/slog"

	"github.com/elect0/chimera/internal/ports"
)

type Service struct {
	log *slog.Logger
}

func NewService(log *slog.Logger) *Service {
	return &Service{
		log: log,
	}
}

func (s *Service) Process(ctx context.Context, opts TransformationOptions, imageBuffer []byte) ([]byte, error) {
	s.log.Info("processing image", slog.Int("width", opts.Width), slog.Int("height", opts.Height))
	return imageBuffer, nil
}

var _ ports.TransformationService = (*Service)(nil)
