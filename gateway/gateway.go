package gateway

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/kanthorlabs/common/gateway/config"
	"github.com/kanthorlabs/common/gateway/httpx"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/patterns"
)

type Gateway interface {
	patterns.Runnable
	UseHttpx(handler httpx.Httpx) error
}

func New(conf *config.Config, logger logging.Logger) (Gateway, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	logger = logger.With("gateway", "default")
	return &gw{conf: conf, logger: logger}, nil
}

type gw struct {
	conf   *config.Config
	logger logging.Logger

	mu     sync.Mutex
	status int
	httpx  *http.Server
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
	if gateway.httpx != nil {
		if err := gateway.httpx.Shutdown(ctx); err != nil {
			returning = errors.Join(returning, err)
		}
	}

	return returning
}

func (gateway *gw) Run(ctx context.Context) error {
	if gateway.httpx != nil {
		go gateway.serve()
		return nil
	}

	return ErrHandlerNotSet
}

func (gateway *gw) serve() {
	err := gateway.httpx.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		gateway.logger.Error(err)
	}
}

func (gateway *gw) UseHttpx(handler httpx.Httpx) error {
	gateway.mu.Lock()
	defer gateway.mu.Unlock()

	if gateway.httpx != nil {
		return ErrHandlerAlreadySet
	}

	gateway.httpx = &http.Server{
		Addr:    gateway.conf.Addr,
		Handler: handler,
	}
	return nil
}
