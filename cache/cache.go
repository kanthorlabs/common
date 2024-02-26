package cache

import (
	"context"
	"time"

	"github.com/kanthorlabs/common/patterns"
)

type Cache interface {
	patterns.Connectable
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, entry any, ttl time.Duration) error
	Exist(ctx context.Context, key string) bool
	Del(ctx context.Context, key string) error
	Expire(ctx context.Context, key string, at time.Time) error
}
