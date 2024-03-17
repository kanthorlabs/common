package streaming

import (
	"context"

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
	returning := safe.Map[error]{}
	p := pool.New().
		WithContext(ctx).
		WithMaxGoroutines(publisher.conf.Publisher.RateLimit)

	for id := range events {
		rid := id
		event := events[id]

		if err := event.Validate(); err != nil {
			publisher.logger.Errorw(ErrPubEventValidation.Error(), "event", event.String(), "error", err.Error())
			returning.Set(rid, err)
			continue
		}

		msg := NatsMsgFromEvent(event)
		p.Go(func(subctx context.Context) error {
			// We don't want to return the error inside this goroutine
			// because we want to set error for each id individually, not merge all of them

			// we will let jetstream handle the context timeout by themself
			ack, err := publisher.js.PublishMsg(subctx, msg)
			if err != nil {
				publisher.logger.Errorw(ErrPubEventPublish.Error(), "event", event.String(), "error", err.Error())
				returning.Set(rid, err)
				return nil
			}

			if ack.Duplicate {
				publisher.logger.Errorw(ErrPubEventDuplicated.Error(), "event", event.String())
				returning.Set(rid, ErrPubEventDuplicated)
				return nil
			}

			return nil
		})
	}

	// no error to handle error because we didn't return it in the p.Go
	p.Wait()

	return returning.Data()
}
