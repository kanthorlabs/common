package config

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestStrategy(t *testing.T) {
	t.Run("KO - enginee error", func(st *testing.T) {
		conf := &Strategy{}
		require.ErrorContains(st, conf.Validate(), "PASSPORT.STRATEGY.CONFIG.ENGINE")
	})

	t.Run("KO - Ask error", func(st *testing.T) {
		conf := &Strategy{Engine: EngineAsk, Name: uuid.NewString()}
		require.ErrorContains(st, conf.Validate(), "PASSPORT.STRATEGY.ASK.CONFIG")
	})

	t.Run("KO - Durability error", func(st *testing.T) {
		conf := &Strategy{Engine: EngineDurability, Name: uuid.NewString()}
		require.ErrorContains(st, conf.Validate(), "SQLX.CONFIG.")
	})
}
