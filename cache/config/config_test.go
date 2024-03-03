package config

import (
	"testing"

	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		conf := &Config{
			Uri: testdata.RedisUri,
		}
		require.NoError(t, conf.Validate())
	})

	t.Run("KO - empty", func(st *testing.T) {
		conf := &Config{}
		require.ErrorContains(t, conf.Validate(), "CACHE.CONFIG.")
	})
}
