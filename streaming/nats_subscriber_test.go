package streaming

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	mockjetstream "github.com/kanthorlabs/common/mocks/jetstream"
	"github.com/kanthorlabs/common/project"
	"github.com/kanthorlabs/common/streaming/entities"
	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/testify"
	natsio "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNatsSubscriber_Name(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		js := mockjetstream.NewJetStream(t)
		subscriber := &NatsSubscriber{
			name:   pubsubname(),
			conf:   testconf("nats://127.0.0.1:42222"),
			logger: testify.Logger(),
			js:     js,
		}
		require.Equal(st, subscriber.name, subscriber.Name())
	})
}

func TestNatsSubscriber_Connect(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		js := mockjetstream.NewJetStream(st)
		subscriber := &NatsSubscriber{
			name:   pubsubname(),
			conf:   testconf("nats://127.0.0.1:42222"),
			logger: testify.Logger(),
			js:     js,
		}
		require.NoError(st, subscriber.Connect(context.Background()))
	})

	t.Run("KO - already connected error", func(st *testing.T) {
		js := mockjetstream.NewJetStream(st)
		subscriber := &NatsSubscriber{
			name:   pubsubname(),
			conf:   testconf("nats://127.0.0.1:42222"),
			logger: testify.Logger(),
			js:     js,
		}
		require.NoError(st, subscriber.Connect(context.Background()))
		require.ErrorIs(st, subscriber.Connect(context.Background()), ErrSubAlreadyConnected)
	})
}

func TestNatsSubscriber_Readiness(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		js := mockjetstream.NewJetStream(st)
		subscriber := &NatsSubscriber{
			name:   pubsubname(),
			conf:   testconf("nats://127.0.0.1:42222"),
			logger: testify.Logger(),
			js:     js,
		}
		require.NoError(st, subscriber.Connect(context.Background()))
		js.EXPECT().Stream(mock.Anything, subscriber.conf.Name).Return(nil, nil)
		require.NoError(st, subscriber.Readiness())
	})

	t.Run("OK - disconnected", func(st *testing.T) {
		js := mockjetstream.NewJetStream(st)
		subscriber := &NatsSubscriber{
			name:   pubsubname(),
			conf:   testconf("nats://127.0.0.1:42222"),
			logger: testify.Logger(),
			js:     js,
		}
		require.NoError(st, subscriber.Connect(context.Background()))
		require.NoError(st, subscriber.Disconnect(context.Background()))
		require.NoError(st, subscriber.Readiness())
	})

	t.Run("KO - not connected error", func(st *testing.T) {
		js := mockjetstream.NewJetStream(st)
		subscriber := &NatsSubscriber{
			name:   pubsubname(),
			conf:   testconf("nats://127.0.0.1:42222"),
			logger: testify.Logger(),
			js:     js,
		}
		require.ErrorIs(st, subscriber.Readiness(), ErrSubNotConnected)
	})
}

func TestNatsSubscriber_Liveness(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		js := mockjetstream.NewJetStream(st)
		subscriber := &NatsSubscriber{
			name:   pubsubname(),
			conf:   testconf("nats://127.0.0.1:42222"),
			logger: testify.Logger(),
			js:     js,
		}
		require.NoError(st, subscriber.Connect(context.Background()))
		js.EXPECT().Stream(mock.Anything, subscriber.conf.Name).Return(nil, nil)
		require.NoError(st, subscriber.Liveness())
	})

	t.Run("OK - disconnected", func(st *testing.T) {
		js := mockjetstream.NewJetStream(st)
		subscriber := &NatsSubscriber{
			name:   pubsubname(),
			conf:   testconf("nats://127.0.0.1:42222"),
			logger: testify.Logger(),
			js:     js,
		}
		require.NoError(st, subscriber.Connect(context.Background()))
		require.NoError(st, subscriber.Disconnect(context.Background()))
		require.NoError(st, subscriber.Liveness())
	})

	t.Run("KO - not connected error", func(st *testing.T) {
		js := mockjetstream.NewJetStream(st)
		subscriber := &NatsSubscriber{
			name:   pubsubname(),
			conf:   testconf("nats://127.0.0.1:42222"),
			logger: testify.Logger(),
			js:     js,
		}
		require.ErrorIs(st, subscriber.Liveness(), ErrSubNotConnected)
	})
}

func TestNatsSubscriber_Disconnect(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		js := mockjetstream.NewJetStream(st)
		subscriber := &NatsSubscriber{
			name:   pubsubname(),
			conf:   testconf("nats://127.0.0.1:42222"),
			logger: testify.Logger(),
			js:     js,
		}
		require.NoError(st, subscriber.Connect(context.Background()))
		require.NoError(st, subscriber.Disconnect(context.Background()))
	})

	t.Run("KO - already connected error", func(st *testing.T) {
		js := mockjetstream.NewJetStream(st)
		subscriber := &NatsSubscriber{
			name:   pubsubname(),
			conf:   testconf("nats://127.0.0.1:42222"),
			logger: testify.Logger(),
			js:     js,
		}
		require.ErrorIs(st, subscriber.Disconnect(context.Background()), ErrSubNotConnected)
	})
}

func TestNatsSubscriber_Subscribe(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		js := mockjetstream.NewJetStream(st)
		subscriber := &NatsSubscriber{
			name:   pubsubname(),
			conf:   testconf("nats://127.0.0.1:42222"),
			logger: testify.Logger(),
			js:     js,
		}
		require.NoError(st, subscriber.Connect(context.Background()))

		batch := mockjetstream.NewMessageBatch(st)
		ch := make(chan jetstream.Msg, 1)
		batch.EXPECT().Messages().Return(ch)
		batch.EXPECT().Error().Return(testdata.ErrGeneric)

		consumer := mockjetstream.NewConsumer(st)
		consumer.EXPECT().CachedInfo().Return(&jetstream.ConsumerInfo{})
		consumer.EXPECT().Fetch(
			subscriber.conf.Subscriber.Concurrency,
			mock.AnythingOfType("jetstream.FetchOpt"),
		).
			Return(batch, nil).
			After(time.Millisecond * time.Duration(subscriber.conf.Subscriber.Timeout))

		topic := topicname()
		js.EXPECT().
			CreateOrUpdateConsumer(mock.Anything, subscriber.conf.Name, consumerconf(subscriber, topic)).
			Return(consumer, nil)

		subject := subjectname()
		ack := &entities.Event{
			Subject: subject,
			Id:      uuid.NewString(),
			Data:    []byte("ack"),
			Metadata: map[string]string{
				"User-Agent": testdata.Fake.UserAgent().InternetExplorer(),
			},
		}
		nak := &entities.Event{
			Subject: subject,
			Id:      uuid.NewString(),
			Data:    []byte("nak"),
			Metadata: map[string]string{
				"User-Agent": testdata.Fake.UserAgent().InternetExplorer(),
			},
		}

		err := subscriber.Sub(context.Background(), topic, consumerhandler(map[string]int{nak.Id: -1, ack.Id: 7}))
		require.NoError(st, err)

		// we expect to receive 2 events but both ack and nak will throw an error
		consumerfetch(st, ch, []*entities.Event{ack, nak}, map[string]int{nak.Id: -1, ack.Id: 5})
		// must close the channel to let batch fetching to finish
		close(ch)

		time.Sleep(time.Second)
		require.NoError(st, subscriber.Disconnect(context.Background()))
		time.Sleep(time.Second)
	})

	t.Run("OK - no valid event", func(st *testing.T) {
		js := mockjetstream.NewJetStream(st)
		subscriber := &NatsSubscriber{
			name:   pubsubname(),
			conf:   testconf("nats://127.0.0.1:42222"),
			logger: testify.Logger(),
			js:     js,
		}
		require.NoError(st, subscriber.Connect(context.Background()))

		batch := mockjetstream.NewMessageBatch(st)
		ch := make(chan jetstream.Msg, 1)
		batch.EXPECT().Messages().Return(ch)

		consumer := mockjetstream.NewConsumer(st)
		consumer.EXPECT().CachedInfo().Return(&jetstream.ConsumerInfo{})
		consumer.EXPECT().Fetch(
			subscriber.conf.Subscriber.Concurrency,
			mock.AnythingOfType("jetstream.FetchOpt"),
		).
			Return(batch, nil).
			After(time.Millisecond * time.Duration(subscriber.conf.Subscriber.Timeout))

		topic := topicname()
		js.EXPECT().
			CreateOrUpdateConsumer(mock.Anything, subscriber.conf.Name, consumerconf(subscriber, topic)).
			Return(consumer, nil)

		err := subscriber.Sub(context.Background(), topic, consumerhandler(nil))
		require.NoError(st, err)

		consumerfetch(st, ch, []*entities.Event{
			{Id: uuid.NewString()},
			{Id: uuid.NewString()},
			{Id: uuid.NewString()},
		}, nil)
		// must close the channel to let batch fetching to finish
		close(ch)

		time.Sleep(time.Second)
		require.NoError(st, subscriber.Disconnect(context.Background()))
		time.Sleep(time.Second)
	})

	t.Run("OK - batch fetching got error", func(st *testing.T) {
		js := mockjetstream.NewJetStream(st)
		subscriber := &NatsSubscriber{
			name:   pubsubname(),
			conf:   testconf("nats://127.0.0.1:42222"),
			logger: testify.Logger(),
			js:     js,
		}
		require.NoError(st, subscriber.Connect(context.Background()))

		consumer := mockjetstream.NewConsumer(st)
		consumer.EXPECT().CachedInfo().Return(&jetstream.ConsumerInfo{})
		consumer.EXPECT().Fetch(
			subscriber.conf.Subscriber.Concurrency,
			mock.AnythingOfType("jetstream.FetchOpt"),
		).
			Return(nil, testdata.ErrGeneric).
			After(time.Millisecond * time.Duration(subscriber.conf.Subscriber.Timeout))

		topic := topicname()
		js.EXPECT().
			CreateOrUpdateConsumer(mock.Anything, subscriber.conf.Name, consumerconf(subscriber, topic)).
			Return(consumer, nil)

		err := subscriber.Sub(context.Background(), topic, consumerhandler(nil))
		require.NoError(st, err)

		time.Sleep(time.Second)
		require.NoError(st, subscriber.Disconnect(context.Background()))
		time.Sleep(time.Second)
	})

	t.Run("KO - not connected error", func(st *testing.T) {
		js := mockjetstream.NewJetStream(st)
		subscriber := &NatsSubscriber{
			name:   pubsubname(),
			conf:   testconf("nats://127.0.0.1:42222"),
			logger: testify.Logger(),
			js:     js,
		}

		topic := topicname()
		err := subscriber.Sub(context.Background(), topic, consumerhandler(nil))
		require.ErrorIs(st, err, ErrSubNotConnected)
	})

	t.Run("KO - topic error", func(st *testing.T) {
		js := mockjetstream.NewJetStream(st)
		subscriber := &NatsSubscriber{
			name:   pubsubname(),
			conf:   testconf("nats://127.0.0.1:42222"),
			logger: testify.Logger(),
			js:     js,
		}
		require.NoError(st, subscriber.Connect(context.Background()))

		topic := topicname() + "#" + topicname()
		err := subscriber.Sub(context.Background(), topic, consumerhandler(nil))
		require.ErrorContains(st, err, "STREAMING.SUBSCRIBER.TOPIC")
	})

	t.Run("KO - consumer error", func(st *testing.T) {
		js := mockjetstream.NewJetStream(st)
		subscriber := &NatsSubscriber{
			name:   pubsubname(),
			conf:   testconf("nats://127.0.0.1:42222"),
			logger: testify.Logger(),
			js:     js,
		}
		require.NoError(st, subscriber.Connect(context.Background()))

		topic := topicname()
		js.EXPECT().
			CreateOrUpdateConsumer(mock.Anything, subscriber.conf.Name, consumerconf(subscriber, topic)).
			Return(nil, testdata.ErrGeneric)

		err := subscriber.Sub(context.Background(), topic, consumerhandler(nil))
		require.ErrorIs(st, err, testdata.ErrGeneric)
	})
}

func consumerconf(subscriber *NatsSubscriber, topic string) jetstream.ConsumerConfig {
	return jetstream.ConsumerConfig{
		Name:            subscriber.name,
		FilterSubject:   fmt.Sprintf("%s.>", project.Subject(topic)),
		MaxDeliver:      subscriber.conf.Subscriber.MaxRetry + 1,
		AckWait:         time.Millisecond * time.Duration(subscriber.conf.Subscriber.Timeout),
		MaxRequestBatch: subscriber.conf.Subscriber.Concurrency,
		DeliverPolicy:   jetstream.DeliverAllPolicy,
		AckPolicy:       jetstream.AckExplicitPolicy,
	}
}

func consumerhandler(ck map[string]int) SubHandler {
	return func(ctx context.Context, events map[string]*entities.Event) map[string]error {
		returning := map[string]error{}
		if len(ck) > 0 {
			for id, status := range ck {
				// 001
				if status == -1 {
					returning[id] = testdata.ErrGeneric
				}
				// 101
				if status == 5 {
					returning[id] = testdata.ErrGeneric
				}
			}
		}
		return returning
	}
}

func consumerfetch(t *testing.T, ch chan jetstream.Msg, events []*entities.Event, ck map[string]int) {
	for _, event := range events {
		jsmsg := mockjetstream.NewMsg(t)
		jsmsg.EXPECT().Subject().Return(event.Subject).Times(1)
		jsmsg.EXPECT().Headers().Return(natsio.Header{
			natsio.MsgIdHdr: []string{event.Id},
			"User-Agent":    []string{event.Metadata["User-Agent"]},
		}).Times(3)
		jsmsg.EXPECT().Data().Return(event.Data).Times(1)

		if len(ck) > 0 {
			status := ck[event.Id]

			// 001
			if status == -1 {
				jsmsg.EXPECT().Nak().Return(testdata.ErrGeneric).Times(1)
			}
			// 011
			if status == -3 {
				jsmsg.EXPECT().Nak().Return(nil).Times(1)
			}
			// 101
			if status == 5 {
				jsmsg.EXPECT().Ack().Return(testdata.ErrGeneric).Times(1)
			}
			// 111
			if status == 7 {
				jsmsg.EXPECT().Nak().Return(nil).Times(1)
			}
		}
		ch <- jsmsg
	}
}
