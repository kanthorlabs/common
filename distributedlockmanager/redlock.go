package distributedlockmanager

import (
	"context"
	"errors"
	"time"

	"github.com/go-redsync/redsync/v4"
	wrapper "github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/kanthorlabs/common/distributedlockmanager/config"
	goredis "github.com/redis/go-redis/v9"
)

func NewRedlock(conf *config.Config) (Factory, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	opts, err := goredis.ParseURL(conf.Uri)
	if err != nil {
		return nil, err
	}

	client := goredis.NewClient(opts)
	rs := redsync.New(wrapper.NewPool(client))

	return func(key string, opts ...config.Option) DistributedLockManager {
		k, err := Key(key)
		instance := &redlock{key: k, err: err}

		if instance.err == nil {
			fconf := &config.Config{Uri: conf.Uri, TimeToLive: conf.TimeToLive}
			for _, opt := range opts {
				opt(fconf)
			}
			instance.conf = fconf
			instance.mu = rs.NewMutex(k, redsync.WithExpiry(time.Millisecond*time.Duration(fconf.TimeToLive)))
		}

		return instance
	}, nil
}

type redlock struct {
	key string
	err error

	conf *config.Config
	mu   *redsync.Mutex
}

func (dlm *redlock) Lock(ctx context.Context) error {
	if dlm.err != nil {
		return dlm.err
	}
	if err := dlm.mu.LockContext(ctx); err != nil {
		return errors.New("DISTRIBUTED_LOCK_MANAGER.LOCK.ERROR")
	}
	return nil
}

func (dlm *redlock) Unlock(ctx context.Context) error {
	if dlm.err != nil {
		return dlm.err
	}
	ok, err := dlm.mu.UnlockContext(ctx)
	if err != nil || !ok {
		return errors.New("DISTRIBUTED_LOCK_MANAGER.UNLOCK.ERROR")
	}

	return nil
}
