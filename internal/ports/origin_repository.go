package ports

import "context"

type OriginRepository interface {
	Get (ctx context.Context, imagePath string) ([]byte, error)
}
