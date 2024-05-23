package cache

import (
	"testing"

	"github.com/kanthorlabs/common/cache/config"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestCache_New(t *testing.T) {
	t.Run("OK - redis", func(st *testing.T) {
		conf := &config.Config{
			Uri: testdata.RedisUri,
		}
		_, err := New(conf)
		require.NoError(st, err)
	})

	t.Run("KO - unknown error", func(st *testing.T) {
		conf := &config.Config{
			Uri: "tcp://127.0.0.1",
		}
		_, err := New(conf)
		require.ErrorContains(st, err, "CACHE.SCHEME_UNKNOWN.ERROR")
	})
}
