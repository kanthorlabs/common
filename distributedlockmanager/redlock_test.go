package distributedlockmanager

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/distributedlockmanager/config"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestRedlock(t *testing.T) {
	testconf := &config.Config{
		Uri:        os.Getenv("REDIS_URI"),
		TimeToLive: testdata.Fake.UInt64Between(10000, 100000),
	}
	if testconf.Uri == "" {
		testconf.Uri = "redis://localhost:6379/0"
	}

	t.Run("New", func(st *testing.T) {
		st.Run("KO - configuration error", func(sst *testing.T) {
			conf := &config.Config{}
			_, err := NewRedlock(conf)
			require.ErrorContains(st, err, "DISTRIBUTED_LOCK_MANAGER.CONFIG.")
		})
		st.Run("KO - redis url error", func(sst *testing.T) {
			conf := &config.Config{
				Uri:        "tcp://localhost:6379/0",
				TimeToLive: testdata.Fake.UInt64Between(10000, 100000),
			}
			_, err := NewRedlock(conf)
			require.ErrorContains(st, err, "redis: ")
		})
	})

	t.Run(".Lock/.Unlock", func(st *testing.T) {
		dlm, err := NewRedlock(testconf)
		require.Nil(st, err)

		key := uuid.NewString()
		locker := dlm(key, config.TimeToLive(testdata.Fake.UInt64Between(10000, 100000)))

		err = locker.Lock(context.Background())
		require.Nil(st, err)

		err = locker.Unlock(context.Background())
		require.Nil(st, err)

		err = locker.Unlock(context.Background())
		require.ErrorContains(st, err, "DISTRIBUTED_LOCK_MANAGER.UNLOCK.ERROR")
	})
}
