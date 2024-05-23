package idempotency

import (
	"context"
	"errors"
	"strings"

	"github.com/kanthorlabs/common/idempotency/config"
	"github.com/kanthorlabs/common/patterns"
)

// New creates a new idempotency instance based on the provided configuration.
// The idempotency instance is initialized based on the URI scheme.
// Supported schemes are:
// - memory://
// - redis://
// If the URI scheme is not supported, an error is returned.
func New(conf *config.Config) (Idempotency, error) {
	if strings.HasPrefix(conf.Uri, "redis") {
		return NewRedis(conf)
	}

	return nil, errors.New("IDEMPOTENCY.SCHEME_UNKNOWN.ERROR")
}

type Idempotency interface {
	patterns.Connectable
	Validate(ctx context.Context, key string) error
}
