package distributedlockmanager

import (
	"testing"

	"github.com/kanthorlabs/common/distributedlockmanager/config"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

var testconf = &config.Config{
	Uri:        "memory://",
	TimeToLive: testdata.Fake.UInt64Between(10000, 100000),
}

func TestIdempotency(t *testing.T) {
	t.Run("OK - memory", func(st *testing.T) {
		conf := &config.Config{
			Uri:        "memory://",
			TimeToLive: testdata.Fake.UInt64Between(10000, 100000),
		}
		_, err := New(conf)
		require.NoError(st, err)
	})

	t.Run("OK - redlock", func(st *testing.T) {
		conf := &config.Config{
			Uri:        "redis://localhost:6379/0",
			TimeToLive: testdata.Fake.UInt64Between(10000, 100000),
		}
		_, err := New(conf)
		require.NoError(st, err)
	})

	t.Run("KO - unknown error", func(st *testing.T) {
		conf := &config.Config{
			Uri:        "tcp://127.0.0.1",
			TimeToLive: testdata.Fake.UInt64Between(10000, 100000),
		}
		_, err := New(conf)
		require.ErrorContains(st, err, "DISTRIBUTED_LOCK_MANAGER.SCHEME_UNKNOWN.ERROR")
	})

	t.Run("KO - configuration error", func(st *testing.T) {
		conf := &config.Config{}
		_, err := New(conf)
		require.ErrorContains(st, err, "DISTRIBUTED_LOCK_MANAGER.CONFIG.")
	})
}
