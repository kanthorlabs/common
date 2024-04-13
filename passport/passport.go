package passport

import (
	"context"
	"sync"

	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/passport/config"
	"github.com/kanthorlabs/common/passport/strategies"
	"github.com/kanthorlabs/common/patterns"
	"github.com/sourcegraph/conc/pool"
)

// New creates a new passport instance with all registered strategies.
// That instance is a facade that allows to interact with all strategies at once.
func New(conf *config.Config, logger logging.Logger) (Passport, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	instances := make(map[string]strategies.Strategy)

	for i := range conf.Strategies {
		engine := conf.Strategies[i].Engine
		name := conf.Strategies[i].Name

		if _, exist := instances[name]; exist {
			return nil, ErrStrategyDuplicated
		}

		if engine == config.EngineAsk {
			strategy, err := strategies.NewAsk(
				&conf.Strategies[i].Ask,
				logger.With("strategy_engine", engine, "strategy_name", name),
			)
			if err != nil {
				return nil, err
			}

			instances[name] = strategy
		}

		if engine == config.EngineInternal {
			strategy, err := strategies.NewInternal(
				&conf.Strategies[i].Internal,
				logger.With("strategy_engine", engine, "strategy_name", name),
			)
			if err != nil {
				return nil, err
			}

			instances[conf.Strategies[i].Name] = strategy
		}

		if engine == config.EngineExternal {
			strategy, err := strategies.NewExternal(
				&conf.Strategies[i].External,
				logger.With("strategy_engine", engine, "strategy_name", name),
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
	Strategy(name string) (strategies.Strategy, error)
}

type passport struct {
	strategies map[string]strategies.Strategy

	mu     sync.Mutex
	status int
}

func (instance *passport) Connect(ctx context.Context) error {
	instance.mu.Lock()
	defer instance.mu.Unlock()

	if instance.status == patterns.StatusConnected {
		return ErrAlreadyConnected
	}

	p := pool.New().WithErrors()
	for i := range instance.strategies {
		strategy := instance.strategies[i]
		p.Go(func() error {
			return strategy.Connect(ctx)
		})
	}
	if err := p.Wait(); err != nil {
		return err
	}

	instance.status = patterns.StatusConnected
	return nil
}

func (instance *passport) Readiness() error {
	if instance.status == patterns.StatusDisconnected {
		return nil
	}
	if instance.status != patterns.StatusConnected {
		return ErrNotConnected
	}

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
	if instance.status == patterns.StatusDisconnected {
		return nil
	}
	if instance.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	p := pool.New().WithErrors()
	for i := range instance.strategies {
		strategy := instance.strategies[i]
		p.Go(func() error {
			return strategy.Liveness()
		})
	}
	return p.Wait()
}

func (instance *passport) Disconnect(ctx context.Context) error {
	instance.mu.Lock()
	defer instance.mu.Unlock()

	if instance.status != patterns.StatusConnected {
		return ErrNotConnected
	}
	instance.status = patterns.StatusDisconnected

	p := pool.New().WithErrors()
	for i := range instance.strategies {
		strategy := instance.strategies[i]
		p.Go(func() error {
			return strategy.Disconnect(ctx)
		})
	}
	return p.Wait()
}

func (instance *passport) Strategy(name string) (strategies.Strategy, error) {
	strategy, exist := instance.strategies[name]
	if !exist {
		return nil, ErrStrategyNotFound
	}

	return strategy, nil
}
