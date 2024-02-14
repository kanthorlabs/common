package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run(".Validate/KO", func(st *testing.T) {
		conf := &Config{}
		require.ErrorContains(st, conf.Validate(), "PASSPORT.CONFIG.STRATEGIES")
	})
}
