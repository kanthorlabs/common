package entities

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCredentials_validate(t *testing.T) {
	t.Run("KO - credentials error", func(st *testing.T) {
		credentials := &Credentials{}
		require.ErrorContains(st, credentials.Validate(), "PASSPORT.CREDENTIALS.")
	})
}
