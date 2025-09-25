package ports

import "context"

type CacheRepository interface {
	Get (ctx context.Context, key string) ([]byte, error)
	Set (ctx context.Context, key string, data []byte) error
}
