package entities

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateCredentialsOnLogin(t *testing.T) {
	t.Run("KO - credentials is nil", func(st *testing.T) {
		require.ErrorContains(st, ValidateCredentialsOnLogin(nil), "PASSPORT.CREDENTIALS")
	})
	t.Run("KO - credentials error", func(st *testing.T) {
		require.ErrorContains(st, ValidateCredentialsOnLogin(&Credentials{}), "PASSPORT.CREDENTIALS.")
	})
}
