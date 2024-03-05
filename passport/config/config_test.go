package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("KO - no strategy", func(st *testing.T) {
		conf := &Config{Strategies: make([]Strategy, 0)}
		require.ErrorContains(st, conf.Validate(), "PASSPORT.CONFIG.STRATEGIES")
	})

	t.Run("KO - strategy error", func(st *testing.T) {
		conf := &Config{Strategies: make([]Strategy, 1)}
		require.ErrorContains(st, conf.Validate(), "PASSPORT.STRATEGY.CONFIG")
	})
}
