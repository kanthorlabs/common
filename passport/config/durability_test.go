package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDurability(t *testing.T) {
	t.Run("KO - sqlx error", func(st *testing.T) {
		conf := &Durability{}
		require.ErrorContains(st, conf.Validate(), "SQLX.CONFIG")
	})
}
