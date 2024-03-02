package circuitbreaker

import (
	"errors"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/kanthorlabs/common/circuitbreaker/config"
	"github.com/kanthorlabs/common/logging"
	"github.com/sony/gobreaker"
)

// NewGoBreaker creates a new circuit breaker instance that is using go-breaker as the implementation and the LRU cache for the storage.
// For each command, it will create a new circuit breaker instance and store it in the LRU cache.
// That means, if the command is not used for a while, it will be removed from the cache to free up the memory.
func NewGoBreaker(conf *config.Config, logger logging.Logger) (CircuitBreaker, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	logger = logger.With("circuitbreaker", "gobreaker")
	breakers, err := lru.New[string, *gobreaker.CircuitBreaker](conf.Size)
	if err != nil {
		return nil, err
	}

	return &sonycb{
		conf:     conf,
		logger:   logger,
		breakers: breakers,
	}, nil
}

type sonycb struct {
	conf   *config.Config
	logger logging.Logger

	breakers *lru.Cache[string, *gobreaker.CircuitBreaker]
}

func (cb *sonycb) Do(cmd string, handler Handler, onError ErrorHandler) (any, error) {
	breaker := cb.get(cmd, onError)
	data, err := breaker.Execute(handler)
	// convert error
	if err != nil {
		return nil, cb.error(err)
	}

	return data, nil
}

func (cb *sonycb) get(cmd string, onError ErrorHandler) *gobreaker.CircuitBreaker {
	if breaker, ok := cb.breakers.Get(cmd); ok {
		return breaker
	}

	breaker := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name: cmd,
		// the maximum number of requests allowed to pass through when the CircuitBreaker is half-open
		MaxRequests: cb.conf.Half.PassthroughRequests,
		// the cyclic period of the closed state for CircuitBreaker to clear the internal Counts
		Interval: time.Millisecond * time.Duration(cb.conf.Close.CleanupInterval),
		// the period of the open state, after which the state of CircuitBreaker becomes half-open
		Timeout: time.Millisecond * time.Duration(cb.conf.Open.Duration),
		// if ReadyToTrip returns true, the CircuitBreaker will be placed into the open state.
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			if counts.ConsecutiveFailures >= cb.conf.Open.Conditions.ErrorConsecutive {
				cb.logger.Warnw("CIRCUIT_BREAKER.STAGE_OPEN.CONSECUTIVE_FAILURES", "consecutive", counts.ConsecutiveFailures, "threshold", cb.conf.Open.Conditions.ErrorConsecutive)
				return true
			}

			reached := counts.TotalFailures > cb.conf.Half.PassthroughRequests
			ratio := float32(counts.TotalFailures) / float32(counts.Requests)
			if reached && ratio >= cb.conf.Open.Conditions.ErrorRatio {
				cb.logger.Warnw("CIRCUIT_BREAKER.STAGE_OPEN.ERROR_RATIO", "consecutive", counts.ConsecutiveFailures, "threshold", cb.conf.Open.Conditions.ErrorConsecutive)
				return true
			}

			return false
		},
		OnStateChange: func(cmd string, from gobreaker.State, to gobreaker.State) {
			cb.logger.Warnw("CIRCUIT_BREAKER.STAGE_CHANGE", "cmd", cmd, "from", from.String(), "to", to.String())
		},
		IsSuccessful: func(err error) bool {
			return onError(err) == nil
		},
	})

	cb.breakers.Add(cmd, breaker)
	return breaker
}

func (cb *sonycb) error(err error) error {
	if errors.Is(err, gobreaker.ErrTooManyRequests) {
		return errors.New("CIRCUIT_BREAKER.TOO_MANY_REQUEST.ERROR")
	}
	if errors.Is(err, gobreaker.ErrOpenState) {
		return errors.New("CIRCUIT_BREAKER.STAGE_OPEN.ERROR")
	}
	return err
}
