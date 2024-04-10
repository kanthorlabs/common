package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExternal(t *testing.T) {
	t.Run("KO - uri error", func(st *testing.T) {
		urierr := &External{}
		require.ErrorContains(st, urierr.Validate(), "PASSPORT.CONFIG.EXTERNAL.URI")
		schemeerr := &External{Uri: "grpc://localhost:9180"}
		require.ErrorContains(st, schemeerr.Validate(), "PASSPORT.CONFIG.EXTERNAL.URI")
	})
}
