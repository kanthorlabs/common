package streaming

import (
	"context"

	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/patterns"
	"github.com/kanthorlabs/common/streaming/config"
	"github.com/kanthorlabs/common/streaming/entities"
)

// New returns a new instance of the streaming service which is using NATS as the default
func New(conf *config.Config, logger logging.Logger) (Stream, error) {
	return NewNats(conf, logger)
}

type Stream interface {
	patterns.Connectable
	Publisher(name string) (Publisher, error)
	Subscriber(name string) (Subscriber, error)
}

type Publisher interface {
	Name() string
	Pub(ctx context.Context, events map[string]*entities.Event) map[string]error
}

type Subscriber interface {
	patterns.Connectable
	Name() string
	Sub(ctx context.Context, topic string, handler SubHandler) error
}

type SubHandler func(ctx context.Context, events map[string]*entities.Event) map[string]error
