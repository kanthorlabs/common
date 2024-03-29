package strategies

import (
	"encoding/base64"
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
		hash, err := password.Hash(passwords[i])
		require.NoError(t, err)
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

var (
	user  = testdata.Fake.Internet().Email()
	pass  = uuid.NewString()
	basic = base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))
)
