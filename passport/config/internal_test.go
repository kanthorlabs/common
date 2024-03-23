package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInternal(t *testing.T) {
	t.Run("KO - sqlx error", func(st *testing.T) {
		conf := &Internal{}
		require.ErrorContains(st, conf.Validate(), "SQLX.CONFIG")
	})
}
