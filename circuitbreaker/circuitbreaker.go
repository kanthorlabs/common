package circuitbreaker

import (
	"github.com/kanthorlabs/common/circuitbreaker/config"
	"github.com/kanthorlabs/common/logging"
)

func New(conf *config.Config, logger logging.Logger) (CircuitBreaker, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	return NewGoBreaker(conf, logger)
}

type CircuitBreaker interface {
	Do(cmd string, onHandle Handler, onError ErrorHandler) (any, error)
}

type Handler func() (any, error)

type ErrorHandler func(err error) error
