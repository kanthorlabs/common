package streaming

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/project"
	"github.com/kanthorlabs/common/streaming/config"
	"github.com/kanthorlabs/common/streaming/entities"
	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/require"

	"github.com/nats-io/nats-server/v2/server"
	natstest "github.com/nats-io/nats-server/v2/test"
)

func TestNats_New(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		server := natsserver()
		defer server.Shutdown()

		conf := testconf(server.ClientURL())
		stream, err := NewNats(conf, testify.Logger())
		require.NoError(st, err)
		require.NotNil(st, stream)
	})

	t.Run("KO - validation error", func(st *testing.T) {
		stream, err := NewNats(&config.Config{}, testify.Logger())
		require.Error(st, err)
		require.Nil(st, stream)
	})
}

func TestNats_Connect(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		server := natsserver()
		defer server.Shutdown()

		conf := testconf(server.ClientURL())
		stream, err := NewNats(conf, testify.Logger())
		require.NoError(st, err)
		require.NoError(st, stream.Connect(context.Background()))
	})

	t.Run("KO - already connected error", func(st *testing.T) {
		server := natsserver()
		defer server.Shutdown()

		conf := testconf(server.ClientURL())
		stream, err := NewNats(conf, testify.Logger())
		require.NoError(st, err)
		require.NoError(st, stream.Connect(context.Background()))
		require.Error(st, stream.Connect(context.Background()))
	})

	t.Run("KO - connection error", func(st *testing.T) {
		conf := testconf(testdata.NatsUri)
		stream, err := NewNats(conf, testify.Logger())
		require.NoError(st, err)
		require.ErrorContains(st, stream.Connect(context.Background()), "nats: ")
	})
}

func TestNats_Readiness(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		server := natsserver()
		defer server.Shutdown()

		conf := testconf(server.ClientURL())
		stream, err := NewNats(conf, testify.Logger())
		require.NoError(st, err)
		require.NoError(st, stream.Connect(context.Background()))
		require.NoError(st, stream.Readiness())
	})

	t.Run(testify.CaseOKDisconnected, func(st *testing.T) {
		server := natsserver()
		defer server.Shutdown()

		conf := testconf(server.ClientURL())

		stream, err := NewNats(conf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, stream.Connect(context.Background()))
		require.NoError(st, stream.Disconnect(context.Background()))
		require.NoError(st, stream.Readiness())
	})

	t.Run(testify.CaseKONotConnectedError, func(st *testing.T) {
		conf := testconf(testdata.NatsUri)
		stream, err := NewNats(conf, testify.Logger())
		require.NoError(st, err)
		require.ErrorIs(st, stream.Readiness(), ErrNotConnected)
	})
}

func TestNats_Liveness(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		server := natsserver()
		defer server.Shutdown()

		conf := testconf(server.ClientURL())
		stream, err := NewNats(conf, testify.Logger())
		require.NoError(st, err)
		require.NoError(st, stream.Connect(context.Background()))
		require.NoError(st, stream.Liveness())
	})

	t.Run(testify.CaseOKDisconnected, func(st *testing.T) {
		server := natsserver()
		defer server.Shutdown()

		conf := testconf(server.ClientURL())

		stream, err := NewNats(conf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, stream.Connect(context.Background()))
		require.NoError(st, stream.Disconnect(context.Background()))
		require.NoError(st, stream.Liveness())
	})

	t.Run(testify.CaseKONotConnectedError, func(st *testing.T) {
		conf := testconf(testdata.NatsUri)
		stream, err := NewNats(conf, testify.Logger())
		require.NoError(st, err)
		require.ErrorIs(st, stream.Liveness(), ErrNotConnected)
	})
}

func TestNats_Disconnect(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		server := natsserver()
		defer server.Shutdown()

		conf := testconf(server.ClientURL())
		stream, err := NewNats(conf, testify.Logger())
		require.NoError(st, err)
		require.NoError(st, stream.Connect(context.Background()))
		require.NoError(st, stream.Disconnect(context.Background()))
	})

	t.Run(testify.CaseKONotConnectedError, func(st *testing.T) {
		conf := testconf(testdata.NatsUri)
		stream, err := NewNats(conf, testify.Logger())
		require.NoError(st, err)
		require.ErrorIs(st, stream.Disconnect(context.Background()), ErrNotConnected)
	})
}

func TestNats_Publisher(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		server := natsserver()
		defer server.Shutdown()

		conf := testconf(server.ClientURL())
		stream, err := NewNats(conf, testify.Logger())
		require.NoError(st, err)
		require.NoError(st, stream.Connect(context.Background()))
		defer stream.Disconnect(context.Background())

		publisher, err := stream.Publisher(pubsubname())
		require.NoError(st, err)
		require.NotNil(st, publisher)
	})

	t.Run("OK - already created", func(st *testing.T) {
		server := natsserver()
		defer server.Shutdown()

		conf := testconf(server.ClientURL())
		stream, err := NewNats(conf, testify.Logger())
		require.NoError(st, err)
		require.NoError(st, stream.Connect(context.Background()))
		defer stream.Disconnect(context.Background())

		name := pubsubname()
		_, err = stream.Publisher(name)
		require.NoError(st, err)
		_, err = stream.Publisher(name)
		require.NoError(st, err)

		require.Equal(st, 1, len(stream.(*nats).publishers))
	})

	t.Run(testify.CaseKONotConnectedError, func(st *testing.T) {
		conf := testconf(testdata.NatsUri)
		stream, err := NewNats(conf, testify.Logger())
		require.NoError(st, err)
		_, err = stream.Publisher(streamname())
		require.ErrorIs(st, err, ErrNotConnected)
	})

	t.Run("KO - name error", func(st *testing.T) {
		server := natsserver()
		defer server.Shutdown()

		conf := testconf(server.ClientURL())
		stream, err := NewNats(conf, testify.Logger())
		require.NoError(st, err)
		require.NoError(st, stream.Connect(context.Background()))
		defer stream.Disconnect(context.Background())

		_, err = stream.Publisher("")
		require.ErrorContains(st, err, "STREAMING.PUBLISHER.NAME")
	})
}

func TestNats_Subscriber(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		server := natsserver()
		defer server.Shutdown()

		conf := testconf(server.ClientURL())
		stream, err := NewNats(conf, testify.Logger())
		require.NoError(st, err)
		require.NoError(st, stream.Connect(context.Background()))
		defer stream.Disconnect(context.Background())

		subscriber, err := stream.Subscriber(pubsubname())
		require.NoError(st, err)
		require.NotNil(st, subscriber)
	})

	t.Run("OK - already created", func(st *testing.T) {
		server := natsserver()
		defer server.Shutdown()

		conf := testconf(server.ClientURL())
		stream, err := NewNats(conf, testify.Logger())
		require.NoError(st, err)
		require.NoError(st, stream.Connect(context.Background()))
		defer stream.Disconnect(context.Background())

		name := pubsubname()
		_, err = stream.Subscriber(name)
		require.NoError(st, err)
		_, err = stream.Subscriber(name)
		require.NoError(st, err)

		require.Equal(st, 1, len(stream.(*nats).subscribers))
	})

	t.Run(testify.CaseKONotConnectedError, func(st *testing.T) {
		conf := testconf(testdata.NatsUri)
		stream, err := NewNats(conf, testify.Logger())
		require.NoError(st, err)
		_, err = stream.Subscriber(streamname())
		require.ErrorIs(st, err, ErrNotConnected)
	})

	t.Run("KO - name error", func(st *testing.T) {
		server := natsserver()
		defer server.Shutdown()

		conf := testconf(server.ClientURL())
		stream, err := NewNats(conf, testify.Logger())
		require.NoError(st, err)
		require.NoError(st, stream.Connect(context.Background()))
		defer stream.Disconnect(context.Background())

		_, err = stream.Subscriber("")
		require.ErrorContains(st, err, "STREAMING.SUBSCRIBER.NAME")
	})
}

func TestNats(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		server := natsserver()
		defer server.Shutdown()

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		conf := testconf(server.ClientURL())
		stream, err := NewNats(conf, testify.Logger())
		require.NoError(st, err)
		require.NoError(st, stream.Connect(ctx))
		defer stream.Disconnect(ctx)

		name := strings.ReplaceAll(uuid.NewString(), "-", "")
		subscriber, err := stream.Subscriber(name)
		require.NoError(st, err)

		topic := topicname()
		count := testdata.Fake.IntBetween(conf.Subscriber.Concurrency+1, conf.Subscriber.Concurrency*2-1)
		items := fakeitems(count)
		datac := make(chan *testdata.User, count)

		require.NoError(st, subscriber.Connect(ctx))
		defer subscriber.Disconnect(ctx)
		err = subscriber.Sub(ctx, project.Subject(topic), func(ctx context.Context, events map[string]*entities.Event) map[string]error {
			returning := map[string]error{}
			for id, event := range events {
				var user testdata.User
				if err := json.Unmarshal(event.Data, &user); err != nil {
					returning[id] = err
					continue
				}
				user.Metadata = map[string]string{}
				for k, v := range event.Metadata {
					user.Metadata[k] = v
				}

				datac <- &user
			}
			return returning
		})
		require.NoError(st, err)

		publisher, err := stream.Publisher(name)
		require.NoError(st, err)

		events := map[string]*entities.Event{}
		for id, item := range items {
			events[id] = &entities.Event{
				Subject:  project.Subject(topic, "user"),
				Id:       id,
				Data:     item.Bytes(),
				Metadata: map[string]string{},
			}

			for k, v := range item.Metadata {
				events[id].Metadata[k] = v
			}
		}
		errs := publisher.Pub(ctx, events)
		require.Equal(st, 0, len(errs))

		var received int
		for user := range datac {
			require.Equal(st, items[user.Id], user)
			delete(items, user.Id)
			received++

			if received == count {
				break
			}
		}

		require.Equal(st, 0, len(items))
	})
}

func testconf(uri string) *config.Config {
	return &config.Config{
		Name: streamname(),
		Uri:  uri,
		Nats: config.Nats{
			Replicas: 0,
			Limits: config.NatsLimits{
				Bytes:    16 * 1024 * 1024 * 1024,
				MsgSize:  1 * 1024 * 1024,
				MsgCount: 30000,
				MsgAge:   1 * 24 * 60 * 60,
			},
		},
		Publisher: config.Publisher{RateLimit: testdata.Fake.IntBetween(1, 1000)},
		Subscriber: config.Subscriber{
			Timeout:     testdata.Fake.Int64Between(3000, 5000),
			MaxRetry:    1,
			Concurrency: 100,
		},
	}
}

func natsserver() *server.Server {
	opts := natstest.DefaultTestOptions
	opts.Port = -1
	opts.JetStream = true
	opts.StoreDir = "/tmp/" + uuid.NewString()
	return natstest.RunServer(&opts)
}
