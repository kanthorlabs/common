package strategies

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/cipher/password"
	"github.com/kanthorlabs/common/passport/entities"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) ([]entities.Account, []string) {
	count := testdata.Fake.IntBetween(5, 10)

	accounts := make([]entities.Account, count)
	passwords := make([]string, count)

	for i := 0; i < count; i++ {
		passwords[i] = uuid.NewString()
		hash, err := password.HashString(passwords[i])
		require.Nil(t, err)
		accounts[i] = entities.Account{
			Username:     uuid.NewString(),
			PasswordHash: hash,
			Name:         testdata.Fake.Internet().User(),
			CreatedAt:    time.Now().UnixMilli(),
			UpdatedAt:    time.Now().UnixMilli(),
		}
	}

	return accounts, passwords
}
