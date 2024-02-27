package idempotency

import (
	"context"
	"errors"
	"strings"

	"github.com/kanthorlabs/common/idempotency/config"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/patterns"
)

func New(conf *config.Config, logger logging.Logger) (Idempotency, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	if strings.HasPrefix(conf.Uri, "memory") {
		return NewMemory(conf, logger)
	}

	if strings.HasPrefix(conf.Uri, "redis") {
		return NewRedis(conf, logger)
	}

	return nil, errors.New("IDEMPOTENCY.SCHEME_UNKNOWN.ERROR")
}

type Idempotency interface {
	patterns.Connectable
	Validate(ctx context.Context, key string) error
}
