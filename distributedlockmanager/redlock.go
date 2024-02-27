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
		k := Key(key)

		cconf := &config.Config{Uri: conf.Uri, TimeToLive: conf.TimeToLive}
		for _, opt := range opts {
			opt(cconf)
		}

		return &redlock{
			key:  k,
			conf: cconf,
			mu:   rs.NewMutex(k, redsync.WithExpiry(time.Millisecond*time.Duration(cconf.TimeToLive))),
		}
	}, nil
}

type redlock struct {
	key string

	conf *config.Config
	mu   *redsync.Mutex
}

func (dlm *redlock) Lock(ctx context.Context) error {
	return dlm.mu.LockContext(ctx)
}

func (dlm *redlock) Unlock(ctx context.Context) error {
	ok, err := dlm.mu.UnlockContext(ctx)
	if err != nil || !ok {
		return errors.New("DISTRIBUTED_LOCK_MANAGER.UNLOCK.ERROR")
	}

	return nil
}
