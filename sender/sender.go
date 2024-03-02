package sender

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/sender/config"
	"github.com/kanthorlabs/common/sender/entities"
)

func New(conf *config.Config, logger logging.Logger) (Send, error) {
	http, err := Http(conf, logger)
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context, r *entities.Request) (*entities.Response, error) {
		uri, err := url.ParseRequestURI(r.Uri)
		if err != nil {
			return nil, errors.New("SENDER.URL.PARSE.ERROR")
		}

		// http & https
		if strings.HasPrefix(uri.Scheme, "http") {
			return http(ctx, r)
		}

		return nil, fmt.Errorf("SENDER.SCHEME.NOT_SUPPORT.ERROR: %s", strings.ToUpper(uri.Scheme))
	}, nil
}

type Send func(context.Context, *entities.Request) (*entities.Response, error)
