package distributedlockmanager

import (
	"testing"

	"github.com/kanthorlabs/common/distributedlockmanager/config"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestDistributedLockManager_New(t *testing.T) {
	t.Run("OK - redlock", func(st *testing.T) {
		conf := &config.Config{
			Uri:        testdata.RedisUri,
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
}
