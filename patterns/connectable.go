package patterns

import "context"

// Connectable define how our external services should be implement
// one of important implementation is healthcheck
// we don't want to start an application that cannot connect to other services like database, cache, ...
type Connectable interface {
	Readiness() error
	Liveness() error

	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
}
