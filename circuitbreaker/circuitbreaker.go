package circuitbreaker

import (
	"github.com/kanthorlabs/common/circuitbreaker/config"
	"github.com/kanthorlabs/common/logging"
)

// New creates a new CircuitBreaker instance that is using go-breaker of Sony by default.
func New(conf *config.Config, logger logging.Logger) (CircuitBreaker, error) {
	return NewGoBreaker(conf, logger)
}

type CircuitBreaker interface {
	Do(cmd string, onHandle Handler, onError ErrorHandler) (any, error)
}

type Handler func() (any, error)

type ErrorHandler func(err error) error
