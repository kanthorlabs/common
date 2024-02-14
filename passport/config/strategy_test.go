package config

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestStrategy(t *testing.T) {
	t.Run(".Validate/KO", func(st *testing.T) {
		conf := &Strategy{}
		require.ErrorContains(st, conf.Validate(), "PASSPORT.STRATEGY.CONFIG.ENGINE")
	})

	t.Run(".Validate/KO - Ask", func(st *testing.T) {
		conf := &Strategy{Engine: EngineAsk, Name: uuid.NewString()}
		require.ErrorContains(st, conf.Validate(), "PASSPORT.STRATEGY.ASK.CONFIG")
	})
}
