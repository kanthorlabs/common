package patterns

import "context"

var (
	StatusDisconnected = -100
	StatusConnected    = 100
)

// Connectable interface is a contract for all external services that our application will connect to
type Connectable interface {
	Connect(ctx context.Context) error

	Readiness() error
	Liveness() error

	Disconnect(ctx context.Context) error
}
