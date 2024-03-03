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

func TestAccount_TableName(t *testing.T) {
	acc := &Account{}
	require.Equal(t, acc.TableName(), project.Name("passport_account"))
}

func TestAccount_Validate(t *testing.T) {
	acc := &Account{}
	require.ErrorContains(t, acc.Validate(), "PASSPORT.ACCOUNT")
}

func TestAccount_Censor(t *testing.T) {
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
	require.Empty(t, censored.PasswordHash)
	require.NotEqual(t, acc.PasswordHash, censored.PasswordHash)

	require.Equal(t, acc.Username, censored.Username)
	require.Equal(t, acc.Metadata.String(), censored.Metadata.String())
	require.Equal(t, acc.CreatedAt, censored.CreatedAt)
	require.Equal(t, acc.UpdatedAt, censored.UpdatedAt)
}

func TestAccount_Active(t *testing.T) {
	acc := &Account{
		Username:     uuid.NewString(),
		PasswordHash: testdata.Fake.Internet().Password(),
		Metadata:     &safe.Metadata{},
		CreatedAt:    time.Now().UnixMilli(),
		UpdatedAt:    time.Now().UnixMilli(),
	}

	activated := clock.NewClock(t)
	// not set the deactivated time
	require.True(t, acc.Active(activated))

	acc.DeactivatedAt = time.Now().UnixMilli()

	activated.EXPECT().Now().Return(time.Now().Add(-time.Hour)).Once()
	// not pass the the deactivated time
	require.True(t, acc.Active(activated))

	activated.EXPECT().Now().Return(time.Now().Add(time.Hour)).Once()
	require.False(t, acc.Active(activated))
}
