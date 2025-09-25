package ports

import (
	"context"

	"github.com/elect0/chimera/internal/domain"
)

type TransformationService interface {
	Process(ctx context.Context, opts domain.TransformationOptions, imagePath string) ([]byte, error)
}
