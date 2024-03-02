package cache

import (
	"testing"

	"github.com/kanthorlabs/common/cache/config"
	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/require"
)

var testconf = &config.Config{
	Uri: "memory://",
}

func TestCache(t *testing.T) {
	t.Run("OK - memory", func(st *testing.T) {
		conf := &config.Config{
			Uri: "memory://",
		}
		_, err := New(conf, testify.Logger())
		require.NoError(st, err)
	})

	t.Run("OK - redis", func(st *testing.T) {
		conf := &config.Config{
			Uri: "redis://localhost:6379/0",
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

	t.Run("KO - configuration error", func(st *testing.T) {
		conf := &config.Config{}
		_, err := New(conf, testify.Logger())
		require.ErrorContains(st, err, "CACHE.CONFIG.")
	})
}
