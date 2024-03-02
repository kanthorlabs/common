package distributedlockmanager

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/distributedlockmanager/config"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestMemory(t *testing.T) {
	t.Run("New", func(st *testing.T) {
		st.Run("KO - configuration error", func(sst *testing.T) {
			conf := &config.Config{}
			_, err := NewMemory(conf)
			require.ErrorContains(st, err, "DISTRIBUTED_LOCK_MANAGER.CONFIG.")
		})
	})

	t.Run(".Lock/.Unlock", func(st *testing.T) {
		dlm, err := NewMemory(testconf)
		require.NoError(st, err)

		key := uuid.NewString()
		locker := dlm(key, config.TimeToLive(testdata.Fake.UInt64Between(10000, 100000)))

		err = locker.Lock(context.Background())
		require.NoError(st, err)

		err = locker.Lock(context.Background())
		require.ErrorContains(st, err, "DISTRIBUTED_LOCK_MANAGER.LOCK.ERROR")

		err = locker.Unlock(context.Background())
		require.NoError(st, err)

		err = locker.Unlock(context.Background())
		require.ErrorContains(st, err, "DISTRIBUTED_LOCK_MANAGER.UNLOCK.ERROR")
	})
}
