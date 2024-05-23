package cache

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/kanthorlabs/common/cache/config"
	"github.com/kanthorlabs/common/patterns"
)

func New(conf *config.Config) (Cache, error) {
	if strings.HasPrefix(conf.Uri, "redis") {
		return NewRedis(conf)
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
