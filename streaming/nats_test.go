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

func TestNats(t *testing.T) {
	t.Run("OK - not enough events", func(st *testing.T) {
		server := natsserver()
		defer server.Shutdown()

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		conf := testconf(server.ClientURL())
		stream, err := NewNats(conf, testify.Logger())
		require.Nil(st, err)
		require.Nil(st, stream.Connect(ctx))
		defer stream.Disconnect(ctx)

		name := strings.ReplaceAll(uuid.NewString(), "-", "")
		subscriber, err := stream.Subscriber(name)
		require.Nil(st, err)

		topic := strings.ReplaceAll(uuid.NewString(), "-", "") + "." + strings.ReplaceAll(uuid.NewString(), "-", "")
		count := testdata.Fake.IntBetween(conf.Subscriber.Concurrency+1, conf.Subscriber.Concurrency*2-1)
		items := fakeitems(count)
		datac := make(chan *testdata.User, count)

		require.Nil(st, subscriber.Connect(ctx))
		defer subscriber.Disconnect(ctx)
		err = subscriber.Sub(ctx, topic, func(ctx context.Context, events map[string]*entities.Event) map[string]error {
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
		require.Nil(st, err)

		publisher, err := stream.Publisher(name)
		require.Nil(st, err)

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
