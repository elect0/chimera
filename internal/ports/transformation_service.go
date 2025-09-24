package ports

import (
	"context"

	"github.com/elect0/chimera/internal/core/transformation"
)

type TransformationService interface {
	Process(ctx context.Context, opts transformation.TransformationOptions, imageBuffer []byte) ([]byte, error)
}
