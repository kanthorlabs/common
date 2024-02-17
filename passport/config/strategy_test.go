package config

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestStrategy(t *testing.T) {
	t.Run(".Validate", func(st *testing.T) {
		st.Run("KO - enginee error", func(sst *testing.T) {
			conf := &Strategy{}
			require.ErrorContains(st, conf.Validate(), "PASSPORT.STRATEGY.CONFIG.ENGINE")
		})

		st.Run("KO - Ask error", func(sst *testing.T) {
			conf := &Strategy{Engine: EngineAsk, Name: uuid.NewString()}
			require.ErrorContains(st, conf.Validate(), "PASSPORT.STRATEGY.ASK.CONFIG")
		})

		st.Run("KO - Durability error", func(sst *testing.T) {
			conf := &Strategy{Engine: EngineDurability, Name: uuid.NewString()}
			require.ErrorContains(st, conf.Validate(), "SQLX.CONFIG.")
		})
	})
}
