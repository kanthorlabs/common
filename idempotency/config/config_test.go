package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("KO", func(st *testing.T) {
		conf := &Config{}
		require.ErrorContains(st, conf.Validate(), "IDEMPOTENCY.CONFIG.")
	})
}
