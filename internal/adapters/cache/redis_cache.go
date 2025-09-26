package cache

import (
	"context"
	"log/slog"
	"time"

	"github.com/elect0/chimera/internal/config"
	"github.com/elect0/chimera/internal/ports"
	"github.com/redis/go-redis/v9"
)

type RedisCacheRepository struct {
	client *redis.Client
	log    *slog.Logger
	ttl    time.Duration
}

func NewRedisCacheRepository(ctx context.Context, cfg *config.Config, log *slog.Logger) (*RedisCacheRepository, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	ttl := 1 * time.Hour

	return &RedisCacheRepository{
		client: client,
		log:    log,
		ttl:    ttl,
	}, nil
}

func (r *RedisCacheRepository) Get(ctx context.Context, key string) ([]byte, error) {
	return r.client.Get(ctx, key).Bytes()
}

func (r *RedisCacheRepository) Set(ctx context.Context, key string, data []byte) error {
	return r.client.Set(ctx, key, data, r.ttl).Err()
}

var _ ports.CacheRepository = (*RedisCacheRepository)(nil)

