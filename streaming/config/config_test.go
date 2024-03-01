package config

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		conf := &Config{
			Name:      streamname(),
			Uri:       "nats://localhost:4222",
			Nats:      *natsconf,
			Publisher: Publisher{RateLimit: testdata.Fake.IntBetween(1, 1000)},
			Subscriber: Subscriber{
				Timeout:     testdata.Fake.Int64Between(1000, 100000),
				MaxRetry:    1,
				Concurrency: testdata.Fake.IntBetween(1, 1000),
			},
		}
		require.Nil(st, conf.Validate())
	})

	t.Run("KO", func(st *testing.T) {
		conf := &Config{
			Name: streamname(),
			Uri:  "http://localhost:4222",
		}
		require.ErrorContains(st, conf.Validate(), "STREAMING.CONFIG.URI")
	})

	t.Run("KO - nats error", func(st *testing.T) {
		conf := &Config{
			Name: streamname(),
			Uri:  "nats://localhost:4222",
			Nats: Nats{},
		}

		require.ErrorContains(st, conf.Validate(), "STREAMING.CONFIG.NATS")
	})

	t.Run("KO - publisher error", func(st *testing.T) {
		conf := &Config{
			Name:      streamname(),
			Uri:       "nats://localhost:4222",
			Nats:      *natsconf,
			Publisher: Publisher{},
		}

		require.ErrorContains(st, conf.Validate(), "STREAMING.CONFIG.PUBLISHER")
	})

	t.Run("KO - subscriber error", func(st *testing.T) {
		conf := &Config{
			Name:      streamname(),
			Uri:       "nats://localhost:4222",
			Nats:      *natsconf,
			Publisher: Publisher{RateLimit: testdata.Fake.IntBetween(1, 1000)},
		}

		require.ErrorContains(st, conf.Validate(), "STREAMING.CONFIG.SUBSCRIBER")
	})
}

func streamname() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "") + "_" + strings.ReplaceAll(uuid.NewString(), "-", "")
}
