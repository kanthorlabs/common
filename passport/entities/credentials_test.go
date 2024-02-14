package entities

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateCredentialsOnLogin(t *testing.T) {
	t.Run("KO", func(st *testing.T) {
		require.ErrorContains(st, ValidateCredentialsOnLogin(&Credentials{}), "PASSPORT.CREDENTIALS")
	})
}
