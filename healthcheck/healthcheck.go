package healthcheck

import (
	"context"
)

type Server interface {
	Connect(ctx context.Context) error
	Readiness(check func() error) error
	Liveness(check func() error) error
	Disconnect(ctx context.Context) error
}

type Client interface {
	Readiness() error
	Liveness() error
}
