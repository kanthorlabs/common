package config

import (
	"testing"

	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

var natsconf = &Nats{
	Replicas: testdata.Fake.IntBetween(0, 100),
	Limits: NatsLimits{
		Bytes:    16 * 1024 * 1024 * 1024,
		MsgSize:  1 * 1024 * 1024,
		MsgCount: 30000,
		MsgAge:   1 * 24 * 60 * 60,
	},
}

func TestNats(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		require.NoError(st, natsconf.Validate())
	})

	t.Run("KO", func(st *testing.T) {
		conf := &Nats{Replicas: -1}
		require.ErrorContains(st, conf.Validate(), "STREAMING.CONFIG.NATS")
	})

	t.Run("KO - limits error", func(st *testing.T) {
		conf := &Nats{
			Replicas: testdata.Fake.IntBetween(0, 100),
		}
		require.ErrorContains(st, conf.Validate(), "STREAMING.CONFIG.NATS.LIMITS")
	})
}
