package cache

import (
	"testing"

	"github.com/kanthorlabs/common/cache/config"
	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/require"
)

var testconf = &config.Config{
	Uri: testdata.MemoryUri,
}

func TestCache_New(t *testing.T) {
	t.Run("OK - memory", func(st *testing.T) {
		conf := &config.Config{
			Uri: testdata.MemoryUri,
		}
		_, err := New(conf, testify.Logger())
		require.NoError(st, err)
	})

	t.Run("OK - redis", func(st *testing.T) {
		conf := &config.Config{
			Uri: testdata.RedisUri,
		}
		_, err := New(conf, testify.Logger())
		require.NoError(st, err)
	})

	t.Run("KO - unknown error", func(st *testing.T) {
		conf := &config.Config{
			Uri: "tcp://127.0.0.1",
		}
		_, err := New(conf, testify.Logger())
		require.ErrorContains(st, err, "CACHE.SCHEME_UNKNOWN.ERROR")
	})
}
