package gateway

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/kanthorlabs/common/gateway/config"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/patterns"
)

type Gateway interface {
	patterns.Runnable
	UseHttpx(handler http.Handler) error
}

func New(conf *config.Config, logger logging.Logger) (Gateway, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	return &gw{conf: conf, logger: logger}, nil
}

type gw struct {
	conf   *config.Config
	logger logging.Logger

	mu     sync.Mutex
	status int
	http   *http.Server
}

func (gateway *gw) Start(ctx context.Context) error {
	gateway.mu.Lock()
	defer gateway.mu.Unlock()

	if gateway.status == patterns.StatusStarted {
		return ErrAlreadyStarted
	}
	gateway.status = patterns.StatusStarted

	return nil
}

func (gateway *gw) Stop(ctx context.Context) error {
	gateway.mu.Lock()
	defer gateway.mu.Unlock()

	if gateway.status != patterns.StatusStarted {
		return ErrNotStarted
	}
	gateway.status = patterns.StatusStopped

	var returning error
	if gateway.http != nil {
		if err := gateway.http.Shutdown(ctx); err != nil {
			returning = errors.Join(returning, err)
		}
	}

	return returning
}

func (gateway *gw) Run(ctx context.Context) error {
	if gateway.http != nil {
		go gateway.serve()
		return nil
	}

	return errors.New("GATEWAY.HANDLER.NOT_SET.ERROR")
}

func (gateway *gw) serve() {
	err := gateway.http.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		gateway.logger.Error(err)
	}
}

func (gateway *gw) UseHttpx(handler http.Handler) error {
	gateway.mu.Lock()
	defer gateway.mu.Unlock()

	if gateway.http != nil {
		return errors.New("GATEWAY.HANDLER.ALREADY_SET.ERROR")
	}

	gateway.http = &http.Server{
		Addr:    gateway.conf.Addr,
		Handler: handler,
	}
	return nil
}
