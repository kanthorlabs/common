package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run(".Validate", func(st *testing.T) {
		st.Run("KO - no strategy", func(sst *testing.T) {
			conf := &Config{Strategies: make([]Strategy, 0)}
			require.ErrorContains(sst, conf.Validate(), "PASSPORT.CONFIG.STRATEGIES")
		})

		st.Run("KO - strategy error", func(sst *testing.T) {
			conf := &Config{Strategies: make([]Strategy, 1)}
			require.ErrorContains(sst, conf.Validate(), "PASSPORT.STRATEGY.CONFIG")
		})
	})
}
