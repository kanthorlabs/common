package idempotency

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/kanthorlabs/common/idempotency/config"
	"github.com/kanthorlabs/common/patterns"
	goredis "github.com/redis/go-redis/v9"
)

// NewRedis creates a new idempotency instance that uses Redis as the underlying storage.
func NewRedis(conf *config.Config) (Idempotency, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	return &redict{conf: conf}, nil
}

type redict struct {
	conf *config.Config

	client *goredis.Client

	mu     sync.Mutex
	status int
}

func (instance *redict) Readiness() error {
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

func (instance *redict) Liveness() error {
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

func (instance *redict) Connect(ctx context.Context) error {
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

func (instance *redict) Disconnect(ctx context.Context) error {
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

func (instance *redict) Validate(ctx context.Context, key string) error {
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
