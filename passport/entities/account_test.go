package entities

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/mocks/clock"
	"github.com/kanthorlabs/common/project"
	"github.com/kanthorlabs/common/safe"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestAccount(t *testing.T) {
	t.Run(".TableName", func(st *testing.T) {
		acc := &Account{}
		require.Equal(st, acc.TableName(), project.Name("passport_account"))
	})

	t.Run(".Validate", func(st *testing.T) {
		acc := &Account{}
		require.ErrorContains(st, acc.Validate(), "PASSPORT.ACCOUNT")
	})

	t.Run(".Censor", func(st *testing.T) {
		acc := &Account{
			Username:     uuid.NewString(),
			PasswordHash: testdata.Fake.Internet().Password(),
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
		require.Equal(st, acc.Metadata.String(), censored.Metadata.String())
		require.Equal(st, acc.CreatedAt, censored.CreatedAt)
		require.Equal(st, acc.UpdatedAt, censored.UpdatedAt)
	})

	t.Run(".Active", func(st *testing.T) {
		acc := &Account{
			Username:     uuid.NewString(),
			PasswordHash: testdata.Fake.Internet().Password(),
			Metadata:     &safe.Metadata{},
			CreatedAt:    time.Now().UnixMilli(),
			UpdatedAt:    time.Now().UnixMilli(),
		}

		activated := clock.NewClock(st)
		// not set the deactivated time
		require.True(st, acc.Active(activated))

		acc.DeactivatedAt = time.Now().UnixMilli()

		activated.EXPECT().Now().Return(time.Now().Add(-time.Hour)).Once()
		// not pass the the deactivated time
		require.True(st, acc.Active(activated))

		activated.EXPECT().Now().Return(time.Now().Add(time.Hour)).Once()
		require.False(st, acc.Active(activated))
	})
}
