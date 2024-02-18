package patterns

import "context"

type Runnable interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Run(ctx context.Context) error
}
