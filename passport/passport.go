package passport

import (
	"context"

	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/passport/config"
	"github.com/kanthorlabs/common/passport/entities"
	"github.com/kanthorlabs/common/passport/strategies"
	"github.com/kanthorlabs/common/patterns"
	"github.com/sourcegraph/conc/pool"
)

func New(conf *config.Config, logger logging.Logger) (Passport, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	instances := make(map[string]strategies.Strategy)

	for i := range conf.Strategies {
		if _, exist := instances[conf.Strategies[i].Name]; exist {
			return nil, ErrStrategyDuplicated
		}

		if conf.Strategies[i].Engine == config.EngineAsk {
			strategy, err := strategies.NewAsk(
				&conf.Strategies[i].Ask,
				logger.With("strategy", config.EngineAsk, "strategy_name", conf.Strategies[i].Name),
			)
			if err != nil {
				return nil, err
			}

			instances[conf.Strategies[i].Name] = strategy
		}
	}

	return &passport{strategies: instances}, nil
}

type Passport interface {
	patterns.Connectable
	Login(ctx context.Context, name string, credentials *entities.Credentials) (*entities.Account, error)
	Logout(ctx context.Context, name string, credentials *entities.Credentials) error
	Verify(ctx context.Context, name string, credentials *entities.Credentials) (*entities.Account, error)
	Register(ctx context.Context, name string, acc *entities.Account) error
}

type passport struct {
	strategies map[string]strategies.Strategy
}

func (instance *passport) Readiness() error {
	p := pool.New().WithErrors()
	for i := range instance.strategies {
		strategy := instance.strategies[i]
		p.Go(func() error {
			return strategy.Readiness()
		})
	}
	return p.Wait()
}

func (instance *passport) Liveness() error {
	p := pool.New().WithErrors()
	for i := range instance.strategies {
		strategy := instance.strategies[i]
		p.Go(func() error {
			return strategy.Liveness()
		})
	}
	return p.Wait()
}

func (instance *passport) Connect(ctx context.Context) error {
	p := pool.New().WithErrors()
	for i := range instance.strategies {
		strategy := instance.strategies[i]
		p.Go(func() error {
			return strategy.Connect(ctx)
		})
	}
	return p.Wait()
}

func (instance *passport) Disconnect(ctx context.Context) error {
	p := pool.New().WithErrors()
	for i := range instance.strategies {
		strategy := instance.strategies[i]
		p.Go(func() error {
			return strategy.Connect(ctx)
		})
	}
	return p.Wait()
}

func (instance *passport) Login(ctx context.Context, name string, credentials *entities.Credentials) (*entities.Account, error) {
	strategy, has := instance.strategies[name]
	if !has {
		return nil, ErrStrategyNotFound
	}

	return strategy.Login(ctx, credentials)
}

func (instance *passport) Logout(ctx context.Context, name string, credentials *entities.Credentials) error {
	strategy, has := instance.strategies[name]
	if !has {
		return ErrStrategyNotFound
	}

	return strategy.Logout(ctx, credentials)
}

func (instance *passport) Verify(ctx context.Context, name string, credentials *entities.Credentials) (*entities.Account, error) {
	strategy, has := instance.strategies[name]
	if !has {
		return nil, ErrStrategyNotFound
	}

	return strategy.Verify(ctx, credentials)
}

func (instance *passport) Register(ctx context.Context, name string, acc *entities.Account) error {
	strategy, has := instance.strategies[name]
	if !has {
		return ErrStrategyNotFound
	}

	return strategy.Register(ctx, acc)
}
