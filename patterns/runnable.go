package patterns

import "context"

var (
	StatusStopped = -200
	StatusStarted = 200
)

// Runnable interface is a contract for all services that our application will run
type Runnable interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Run(ctx context.Context) error
}
