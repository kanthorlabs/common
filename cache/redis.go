package cache

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/kanthorlabs/common/cache/config"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/patterns"
	goredis "github.com/redis/go-redis/v9"
)

func NewRedis(conf *config.Config, logger logging.Logger) (Cache, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	logger = logger.With("cache", "redis")
	return &redis{conf: conf, logger: logger}, nil
}

type redis struct {
	conf   *config.Config
	logger logging.Logger

	client *goredis.Client

	mu     sync.Mutex
	status int
}

func (instance *redis) Readiness() error {
	if instance.status == patterns.StatusDisconnected {
		return nil
	}
	if instance.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return instance.client.Ping(ctx).Err()
}

func (instance *redis) Liveness() error {
	if instance.status == patterns.StatusDisconnected {
		return nil
	}
	if instance.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return instance.client.Ping(ctx).Err()
}

func (instance *redis) Connect(ctx context.Context) error {
	instance.mu.Lock()
	defer instance.mu.Unlock()

	if instance.status == patterns.StatusConnected {
		return ErrAlreadyConnected
	}
	conf, err := goredis.ParseURL(instance.conf.Uri)
	if err != nil {
		return err
	}
	instance.client = goredis.NewClient(conf)

	instance.status = patterns.StatusConnected
	return nil
}

func (instance *redis) Disconnect(ctx context.Context) error {
	instance.mu.Lock()
	defer instance.mu.Unlock()

	if instance.status != patterns.StatusConnected {
		return ErrNotConnected
	}
	instance.status = patterns.StatusDisconnected

	var returning error
	if err := instance.client.Close(); err != nil {
		returning = errors.Join(returning, err)
	}
	instance.client = nil

	return returning
}

func (instance *redis) Get(ctx context.Context, key string) ([]byte, error) {
	entry, err := instance.client.Get(ctx, Key(key)).Bytes()
	// convert error type to detect later
	if errors.Is(err, goredis.Nil) {
		return nil, ErrEntryNotFound
	}

	return entry, err
}

func (instance *redis) Set(ctx context.Context, key string, entry any, ttl time.Duration) error {
	var value []byte
	var err error
	if entry != nil {
		value, err = json.Marshal(entry)
		if err != nil {
			return err
		}
	}
	return instance.client.Set(ctx, Key(key), value, ttl).Err()
}

func (instance *redis) Exist(ctx context.Context, key string) bool {
	entry, err := instance.client.Exists(ctx, Key(key)).Result()
	return err == nil && entry > 0
}

func (instance *redis) Del(ctx context.Context, key string) error {
	return instance.client.Del(ctx, Key(key)).Err()
}

func (instance *redis) Expire(ctx context.Context, key string, at time.Time) error {
	ok, err := instance.client.ExpireAt(ctx, Key(key), at).Result()
	if err != nil {
		return err
	}
	if !ok {
		return ErrEntryNotFound
	}
	return nil
}
