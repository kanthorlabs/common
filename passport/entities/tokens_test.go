package entities

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTokens_validate(t *testing.T) {
	t.Run("KO - credentials error", func(st *testing.T) {
		credentials := &Tokens{}
		require.ErrorContains(st, credentials.Validate(), "PASSPORT.TOKENS.")
	})
}
