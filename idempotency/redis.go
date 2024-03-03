package idempotency

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/kanthorlabs/common/idempotency/config"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/patterns"
	goredis "github.com/redis/go-redis/v9"
)

func NewRedis(conf *config.Config, logger logging.Logger) (Idempotency, error) {
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

func (instance *redis) Validate(ctx context.Context, key string) error {
	k, err := Key(key)
	if err != nil {
		return err
	}

	var incr *goredis.IntCmd
	// While the client sends commands using pipelining,
	// the server will be forced to queue the replies, using memory.
	// So we cannot get the incr.Val inside the pipeline to validate
	_, err = instance.client.Pipelined(ctx, func(pipe goredis.Pipeliner) error {
		incr = pipe.Incr(ctx, k)
		pipe.Expire(ctx, k, time.Millisecond*time.Duration(instance.conf.TimeToLive))
		return nil
	})
	if err != nil || incr.Val() > 1 {
		return ErrConflict
	}
	return nil
}
