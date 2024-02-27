package config

import (
	"testing"

	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("KO", func(st *testing.T) {
		conf := &Config{}
		require.ErrorContains(st, conf.Validate(), "DISTRIBUTED_LOCK_MANAGER.CONFIG.")
	})
}

func TestTimeToLive(t *testing.T) {
	ttl := testdata.Fake.UInt64Between(100000, 1000000)
	conf := &Config{}

	TimeToLive(ttl)(conf)

	require.Equal(t, ttl, conf.TimeToLive)
}
