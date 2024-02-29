package streaming

import (
	"context"
	"errors"

	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/safe"
	"github.com/kanthorlabs/common/streaming/config"
	"github.com/kanthorlabs/common/streaming/entities"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/sourcegraph/conc/pool"
)

type NatsPublisher struct {
	name   string
	conf   *config.Config
	logger logging.Logger

	js jetstream.JetStream
}

func (publisher *NatsPublisher) Name() string {
	return publisher.name
}

func (publisher *NatsPublisher) Pub(ctx context.Context, events map[string]*entities.Event) map[string]error {
	donec := make(chan bool, 1)
	defer close(donec)

	returning := safe.Map[error]{}
	go func() {
		p := pool.New().WithMaxGoroutines(publisher.conf.Publisher.RateLimit)
		for refId, event := range events {
			if err := event.Validate(); err != nil {
				publisher.logger.Errorw("STREAMING.PUBLISHER.EVENT_VALIDATION.ERROR", "event", event.String())
				returning.Set(refId, err)
				continue
			}

			msg := NatsMsgFromEvent(event.Subject, event)
			p.Go(func() {
				ack, err := publisher.js.PublishMsg(ctx, msg)
				if err != nil {
					publisher.logger.Errorw("STREAMING.PUBLISHER.EVENT_PUBLISH.ERROR", "event", event.String())
					returning.Set(refId, err)
					return
				}

				if ack.Duplicate {
					publisher.logger.Errorw("STREAMING.PUBLISHER.EVENT_DUPLICATED.ERROR", "event", event.String())
					returning.Set(refId, errors.New("STREAMING.PUBLISHER.EVENT_DUPLICATED.ERROR"))
					return
				}
			})
		}
		p.Wait()

		donec <- true
	}()

	select {
	case <-donec:
		return returning.Data()
	case <-ctx.Done():
		data := returning.Data()
		for refId := range events {
			if _, exist := data[refId]; !exist {
				data[refId] = ctx.Err()
			}
		}
		return data
	}
}
