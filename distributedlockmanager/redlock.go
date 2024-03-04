package distributedlockmanager

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-redsync/redsync/v4"
	wrapper "github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/kanthorlabs/common/distributedlockmanager/config"
	"github.com/kanthorlabs/common/patterns"
	goredis "github.com/redis/go-redis/v9"
)

// NewRedlock creates a new distributed lock manager instance that uses redis as the underlying storage and redlock as the algorithm.
func NewRedlock(conf *config.Config) (DistributedLockManager, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	return &redlock{conf: conf}, nil
}

type redlock struct {
	conf *config.Config

	gredis *goredis.Client
	client *redsync.Redsync
	mu     sync.Mutex
	status int
}

func (instance *redlock) Connect(ctx context.Context) error {
	instance.mu.Lock()
	defer instance.mu.Unlock()

	if instance.status == patterns.StatusConnected {
		return ErrAlreadyConnected
	}

	opts, err := goredis.ParseURL(instance.conf.Uri)
	if err != nil {
		return err
	}

	instance.gredis = goredis.NewClient(opts)
	instance.client = redsync.New(wrapper.NewPool(instance.gredis))

	instance.status = patterns.StatusConnected
	return nil
}

func (instance *redlock) Readiness() error {
	if instance.status == patterns.StatusDisconnected {
		return nil
	}
	if instance.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return instance.gredis.Ping(ctx).Err()
}

func (instance *redlock) Liveness() error {
	if instance.status == patterns.StatusDisconnected {
		return nil
	}
	if instance.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return instance.gredis.Ping(ctx).Err()
}

func (instance *redlock) Disconnect(ctx context.Context) error {
	instance.mu.Lock()
	defer instance.mu.Unlock()

	if instance.status != patterns.StatusConnected {
		return ErrNotConnected
	}
	instance.status = patterns.StatusDisconnected

	var returning error
	if err := instance.gredis.Close(); err != nil {
		returning = errors.Join(returning, err)
	}
	instance.client = nil

	return returning
}

func (dlm *redlock) Lock(ctx context.Context, key string, opts ...config.Option) (Identifier, error) {
	k, err := Key(key)
	if err != nil {
		return nil, err
	}

	conf := &config.Config{Uri: dlm.conf.Uri, TimeToLive: dlm.conf.TimeToLive}
	for _, opt := range opts {
		opt(conf)
	}

	mu := dlm.client.NewMutex(k, redsync.WithExpiry(time.Millisecond*time.Duration(conf.TimeToLive)))
	if err := mu.LockContext(ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", ErrLock.Error(), err)
	}
	return &ridentifier{k: k, mu: mu}, nil
}

type ridentifier struct {
	k  string
	mu *redsync.Mutex
}

func (identifier *ridentifier) Unlock(ctx context.Context) error {
	ok, err := identifier.mu.UnlockContext(ctx)
	if err != nil || !ok {
		return ErrUnlock
	}

	return nil
}
