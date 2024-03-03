package config

import (
	"testing"

	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		conf := &Config{
			Uri:        testdata.RedisUri,
			TimeToLive: 1000,
		}
		require.NoError(st, conf.Validate())
	})

	t.Run("KO - uri error", func(st *testing.T) {
		conf := &Config{
			Uri:        "invalid",
			TimeToLive: testdata.Fake.UInt64Between(1000, 100000),
		}
		require.ErrorContains(st, conf.Validate(), "DISTRIBUTED_LOCK_MANAGER.CONFIG")
	})
}

func TestTimeToLive(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		conf := &Config{}
		TimeToLive(1000)(conf)
		require.Equal(st, uint64(1000), conf.TimeToLive)
	})
}
