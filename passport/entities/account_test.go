package entities

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/safe"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestAccount(t *testing.T) {
	t.Run(".Validate/KO", func(st *testing.T) {
		acc := &Account{}
		require.ErrorContains(st, acc.Validate(), "PASSPORT.ACCOUNT")
	})

	t.Run(".Censor", func(st *testing.T) {
		acc := &Account{
			Username:     uuid.NewString(),
			PasswordHash: testdata.Fake.Internet().Password(),
			Tenant:       testdata.Fake.App().Name(),
			Metadata:     &safe.Metadata{},
			CreatedAt:    time.Now().UnixMilli(),
			UpdatedAt:    time.Now().UnixMilli(),
		}
		acc.Metadata.Set(uuid.NewString(), testdata.Fake.Internet().URL())

		censored := acc.Censor()

		// Password must be censored
		require.Empty(st, censored.PasswordHash)
		require.NotEqual(st, acc.PasswordHash, censored.PasswordHash)

		require.Equal(st, acc.Username, censored.Username)
		require.Equal(st, acc.Tenant, censored.Tenant)
		require.Equal(st, acc.Metadata.String(), censored.Metadata.String())
		require.Equal(st, acc.CreatedAt, censored.CreatedAt)
		require.Equal(st, acc.UpdatedAt, censored.UpdatedAt)
	})
}
