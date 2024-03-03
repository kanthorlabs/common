package cache

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/kanthorlabs/common/cache/config"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/patterns"
)

// New creates a new cache instance based on the provided configuration.
// The cache instance is initialized based on the URI scheme.
// Supported schemes are:
// - memory://
// - redis://
// If the URI scheme is not supported, an error is returned.
func New(conf *config.Config, logger logging.Logger) (Cache, error) {
	if strings.HasPrefix(conf.Uri, "memory") {
		return NewMemory(conf, logger)
	}

	if strings.HasPrefix(conf.Uri, "redis") {
		return NewRedis(conf, logger)
	}

	return nil, errors.New("CACHE.SCHEME_UNKNOWN.ERROR")
}

type Cache interface {
	patterns.Connectable
	Get(ctx context.Context, key string, entry any) error
	Set(ctx context.Context, key string, entry any, ttl time.Duration) error
	Exist(ctx context.Context, key string) bool
	Del(ctx context.Context, key string) error
	Expire(ctx context.Context, key string, at time.Time) error
}
