package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		conf := Default("test", 1000)
		require.NoError(st, conf.Validate())
	})
	t.Run("KO", func(st *testing.T) {
		conf := &Config{}
		require.ErrorContains(st, conf.Validate(), "HEALTHCHECK.CONFIG")
	})
}
