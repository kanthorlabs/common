package streaming

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/patterns"
	"github.com/kanthorlabs/common/streaming/config"
	"github.com/kanthorlabs/common/streaming/entities"
	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	natsio "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/sourcegraph/conc"
)

type NatsSubscriber struct {
	name   string
	conf   *config.Config
	logger logging.Logger

	js jetstream.JetStream

	mu     sync.Mutex
	status int
}

func (subscriber *NatsSubscriber) Name() string {
	return subscriber.name
}

func (subscriber *NatsSubscriber) Connect(ctx context.Context) error {
	subscriber.mu.Lock()
	defer subscriber.mu.Unlock()

	if subscriber.status == patterns.StatusConnected {
		return ErrSubAlreadyConnected
	}

	subscriber.logger.Info("connected")

	subscriber.status = patterns.StatusConnected
	return nil
}

func (subscriber *NatsSubscriber) Readiness() error {
	if subscriber.status == patterns.StatusDisconnected {
		return nil
	}
	if subscriber.status != patterns.StatusConnected {
		return ErrSubNotConnected
	}

	expires := time.Millisecond * time.Duration(subscriber.conf.Subscriber.Timeout)
	ctx, cancel := context.WithTimeout(context.Background(), expires)
	defer cancel()

	_, err := subscriber.js.Stream(ctx, subscriber.conf.Name)
	return err
}

func (subscriber *NatsSubscriber) Liveness() error {
	if subscriber.status == patterns.StatusDisconnected {
		return nil
	}
	if subscriber.status != patterns.StatusConnected {
		return ErrSubNotConnected
	}

	expires := time.Millisecond * time.Duration(subscriber.conf.Subscriber.Timeout)
	ctx, cancel := context.WithTimeout(context.Background(), expires)
	defer cancel()

	_, err := subscriber.js.Stream(ctx, subscriber.conf.Name)
	return err
}

func (subscriber *NatsSubscriber) Disconnect(ctx context.Context) error {
	subscriber.mu.Lock()
	defer subscriber.mu.Unlock()

	if subscriber.status != patterns.StatusConnected {
		return ErrSubNotConnected
	}
	subscriber.status = patterns.StatusDisconnected
	subscriber.logger.Info("disconnected")

	return nil
}

func (subscriber *NatsSubscriber) Sub(ctx context.Context, topic string, handler SubHandler) error {
	if subscriber.status != patterns.StatusConnected {
		return ErrSubNotConnected
	}

	err := validator.StringAlphaNumericUnderscoreHyphenDot("STREAMING.SUBSCRIBER.TOPIC", topic)()
	if err != nil {
		return err
	}

	consumer, err := subscriber.consumer(ctx, subscriber.name, topic)
	if err != nil {
		return err
	}

	subject := consumer.CachedInfo().Config.FilterSubject
	subscriber.logger.Infow(
		"initialized consumer",
		"consumer_name", consumer.CachedInfo().Name,
		"consumer_created_at", consumer.CachedInfo().Created.Format(time.RFC3339),
		"subject", subject,
	)

	expires := time.Millisecond * time.Duration(subscriber.conf.Subscriber.Timeout)
	go func() {
		for {
			// the subscription is no longer available because we closed it programmatically
			if subscriber.status == patterns.StatusDisconnected {
				return
			}

			handlerctx := context.Background()
			batch, err := consumer.Fetch(
				subscriber.conf.Subscriber.Concurrency,
				jetstream.FetchMaxWait(expires),
			)
			if err != nil {
				subscriber.logger.Errorw(
					"STREAMING.SUBSCRIBER.PULL.ERROR",
					"error", err.Error(),
					"wait_time", fmt.Sprintf("%dms", subscriber.conf.Subscriber.Timeout),
				)
				continue
			}

			messages := map[string]jetstream.Msg{}
			events := make(map[string]*entities.Event)
			for msg := range batch.Messages() {
				eventId := msg.Headers().Get(natsio.MsgIdHdr)
				messages[eventId] = msg

				event := NatsMsgToEvent(msg)
				if err := event.Validate(); err != nil {
					subscriber.logger.Errorw(err.Error(), "nats_msg", utils.Stringify(msg))
					continue
				}

				events[eventId] = event
			}

			if len(events) == 0 {
				continue
			}
			errs := handler(handlerctx, events)

			var wg conc.WaitGroup
			for id := range events {
				event := events[id]

				wg.Go(func() {
					if err, ok := errs[event.Id]; ok && err != nil {
						if err := messages[event.Id].Nak(); err != nil {
							// it's important to log entire event here to trace it in our log system
							subscriber.logger.Errorw(ErrSubNakFail.Error(), "event", event.String())
						}
						return
					}

					if err := messages[event.Id].Ack(); err != nil {
						// it's important to log entire event here to trace it in our log system
						subscriber.logger.Errorw(ErrSubAckFail.Error(), "event", event.String())
					}
				})
			}
			wg.Wait()

			if batch.Error() != nil && !errors.Is(batch.Error(), natsio.ErrTimeout) {
				subscriber.logger.Errorw(
					"STREAMING.SUBSCRIBER.PULL.ERROR",
					"error", batch.Error().Error(),
					"wait_time", fmt.Sprintf("%dms", subscriber.conf.Subscriber.Timeout),
				)
			}
		}
	}()

	subscriber.logger.Infow("subscribed",
		"timeout", subscriber.conf.Subscriber.Timeout,
		"max_retry", subscriber.conf.Subscriber.MaxRetry,
		"concurrency", subscriber.conf.Subscriber.Concurrency,
	)
	return nil
}

func (subscriber *NatsSubscriber) consumer(ctx context.Context, name, topic string) (jetstream.Consumer, error) {
	return subscriber.js.CreateOrUpdateConsumer(ctx, subscriber.conf.Name, jetstream.ConsumerConfig{
		// common config
		Name:          name,
		Durable:       name,
		FilterSubject: fmt.Sprintf("%s.>", topic),

		// advance config
		MaxDeliver: subscriber.conf.Subscriber.MaxRetry + 1,
		AckWait:    time.Millisecond * time.Duration(subscriber.conf.Subscriber.Timeout),
		// if MaxRequestBatch is 1, and we are going to request 2, we will get an error
		MaxRequestBatch: subscriber.conf.Subscriber.Concurrency,
		// internal config
		DeliverPolicy: jetstream.DeliverAllPolicy,
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
}
