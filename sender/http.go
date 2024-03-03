package sender

import (
	"context"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/sender/config"
	"github.com/kanthorlabs/common/sender/entities"
)

func Http(conf *config.Config, logger logging.Logger) (Send, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	client := resty.New().
		SetLogger(logger.With("sender", "http")).
		SetTimeout(time.Millisecond * time.Duration(conf.Timeout)).
		SetRetryCount(conf.Retry.Count).
		SetRetryWaitTime(time.Millisecond * time.Duration(conf.Retry.WaitTime)).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			status := r.StatusCode()
			url := r.Request.URL
			if status >= http.StatusInternalServerError {
				logger.Warnw("SENDER.RETRYING", "status", status, "url", url)
				return true
			}
			return false
		}).
		SetHeaders(conf.Headers)

	return func(ctx context.Context, r *entities.Request) (*entities.Response, error) {
		req := client.R().
			SetContext(ctx).
			SetHeaderMultiValues(r.Headers)

		var err = r.Validate()
		if err != nil {
			return nil, err
		}

		var res *resty.Response
		if r.Method == http.MethodGet {
			res, err = req.Get(r.Uri)
		}
		if r.Method == http.MethodPost {
			res, err = req.SetBody(r.Body).Post(r.Uri)
		}
		if r.Method == http.MethodPut {
			res, err = req.SetBody(r.Body).Put(r.Uri)
		}
		if r.Method == http.MethodPatch {
			res, err = req.SetBody(r.Body).Patch(r.Uri)
		}

		if err != nil {
			return nil, err
		}

		return &entities.Response{
			Status:  res.StatusCode(),
			Headers: res.Header(),
			// follow redirect url and got final url
			// most time the response url is same as request url
			Uri:  res.RawResponse.Request.URL.String(),
			Body: res.Body(),
		}, nil
	}, nil
}
