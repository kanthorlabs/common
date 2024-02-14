package patterns

import "context"

type Runnable interface {
	Stop(ctx context.Context) error
	Start(ctx context.Context) error
	Run(ctx context.Context) error
}
