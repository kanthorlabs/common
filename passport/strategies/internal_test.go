package strategies

import (
	"context"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/cipher/password"
	"github.com/kanthorlabs/common/passport/config"
	"github.com/kanthorlabs/common/passport/entities"
	"github.com/kanthorlabs/common/passport/utils"
	sqlxconfig "github.com/kanthorlabs/common/persistence/sqlx/config"
	"github.com/kanthorlabs/common/safe"
	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/require"
)

func TestInternal_New(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		_, err := NewInternal(internalconf, testify.Logger())
		require.NoError(st, err)
	})

	t.Run("KO - configuration error", func(st *testing.T) {
		_, err := NewInternal(&config.Internal{}, testify.Logger())
		require.ErrorContains(st, err, "SQLX.CONFIG")
	})

	t.Run("KO - sqlx error", func(st *testing.T) {
		conf := &config.Internal{Sqlx: sqlxconfig.Config{}}
		_, err := NewInternal(conf, testify.Logger())
		require.ErrorContains(st, err, "SQLX.CONFIG")
	})
}

func TestInternal_Connect(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		c, err := NewInternal(internalconf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, c.Connect(context.Background()))
	})

	t.Run("KO - already connected error", func(st *testing.T) {
		c, err := NewInternal(internalconf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, c.Connect(context.Background()))
		require.ErrorIs(st, c.Connect(context.Background()), ErrAlreadyConnected)
	})
}

func TestInternal_Readiness(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		c, err := NewInternal(internalconf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, c.Connect(context.Background()))
		require.NoError(st, c.Readiness())
	})

	t.Run("OK - disconnected", func(st *testing.T) {
		c, err := NewInternal(internalconf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, c.Connect(context.Background()))
		require.NoError(st, c.Disconnect(context.Background()))
		require.NoError(st, c.Readiness())
	})

	t.Run("KO - not connected error", func(st *testing.T) {
		c, err := NewInternal(internalconf, testify.Logger())
		require.NoError(st, err)

		require.ErrorIs(st, c.Readiness(), ErrNotConnected)
	})
}

func TestInternal_Liveness(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		c, err := NewInternal(internalconf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, c.Connect(context.Background()))
		require.NoError(st, c.Liveness())
	})

	t.Run("OK - disconnected", func(st *testing.T) {
		c, err := NewInternal(internalconf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, c.Connect(context.Background()))
		require.NoError(st, c.Disconnect(context.Background()))
		require.NoError(st, c.Liveness())
	})

	t.Run("KO - not connected error", func(st *testing.T) {
		c, err := NewInternal(internalconf, testify.Logger())
		require.NoError(st, err)

		require.ErrorIs(st, c.Liveness(), ErrNotConnected)
	})
}

func TestInternal_Disconnect(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		c, err := NewInternal(internalconf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, c.Connect(context.Background()))
		require.NoError(st, c.Disconnect(context.Background()))
	})

	t.Run("KO - not connected error", func(st *testing.T) {
		c, err := NewInternal(internalconf, testify.Logger())
		require.NoError(st, err)

		require.ErrorIs(st, c.Disconnect(context.Background()), ErrNotConnected)
	})
}

func TestInternal_Register(t *testing.T) {
	accounts, _ := setup(t)

	conf := &config.Internal{Sqlx: sqlxconfig.Config{
		Uri: testdata.SqliteUri,
		Connection: sqlxconfig.Connection{
			MaxLifetime:  sqlxconfig.DefaultConnMaxLifetime,
			MaxIdletime:  sqlxconfig.DefaultConnMaxIdletime,
			MaxIdleCount: sqlxconfig.DefaultConnMaxIdleCount,
			MaxOpenCount: sqlxconfig.DefaultConnMaxOpenCount,
		},
	}}
	strategy, err := NewInternal(conf, testify.Logger())
	require.NoError(t, err)

	strategy.Connect(context.Background())
	defer strategy.Disconnect(context.Background())

	orm := strategy.(*internal).orm
	require.NoError(t, orm.Create(accounts).Error)

	t.Run("OK", func(st *testing.T) {
		pass := uuid.NewString()
		hash, err := password.Hash(pass)
		require.NoError(st, err)

		acc := entities.Account{
			Username:     uuid.NewString(),
			PasswordHash: hash,
			Name:         testdata.Fake.Internet().User(),
			CreatedAt:    time.Now().UnixMilli(),
			UpdatedAt:    time.Now().UnixMilli(),
		}

		require.NoError(st, strategy.Register(context.Background(), acc))
	})

	t.Run("KO - already exist error", func(st *testing.T) {
		i := testdata.Fake.IntBetween(0, len(accounts)-1)
		err := strategy.Register(context.Background(), accounts[i])
		require.ErrorIs(st, err, ErrRegister)
	})
}

func TestInternal_Login(t *testing.T) {
	accounts, passwords := setup(t)

	strategy, err := NewInternal(internalconf, testify.Logger())
	require.NoError(t, err)

	strategy.Connect(context.Background())
	defer strategy.Disconnect(context.Background())

	orm := strategy.(*internal).orm
	require.NoError(t, orm.Create(accounts).Error)

	t.Run("KO - unimplement error", func(st *testing.T) {
		i := testdata.Fake.IntBetween(0, len(passwords)-1)
		credentials := entities.Credentials{
			Username: accounts[i].Username,
			Password: passwords[i],
		}
		_, err := strategy.Login(context.Background(), credentials)
		require.ErrorContains(st, err, "UNIMPLEMENT.ERROR")
	})
}

func TestInternal_Logout(t *testing.T) {
	accounts, _ := setup(t)

	strategy, err := NewInternal(internalconf, testify.Logger())
	require.NoError(t, err)

	strategy.Connect(context.Background())
	defer strategy.Disconnect(context.Background())

	orm := strategy.(*internal).orm
	require.NoError(t, orm.Create(accounts).Error)

	t.Run("KO - unimplement error", func(st *testing.T) {
		err := strategy.Logout(context.Background(), entities.Tokens{})
		require.ErrorContains(st, err, "UNIMPLEMENT.ERROR")
	})
}

func TestInternal_Verify(t *testing.T) {
	accounts, passwords := setup(t)

	strategy, err := NewInternal(internalconf, testify.Logger())
	require.NoError(t, err)

	strategy.Connect(context.Background())
	defer strategy.Disconnect(context.Background())

	orm := strategy.(*internal).orm
	require.NoError(t, orm.Create(accounts).Error)

	t.Run("OK", func(st *testing.T) {
		i := testdata.Fake.IntBetween(0, len(passwords)-1)
		user := accounts[i].Username
		pass := passwords[i]

		basic := utils.CreateRegionalBasicCredentials(user + ":" + pass)
		tokens := entities.Tokens{Access: utils.SchemeBasic + basic}
		acc, err := strategy.Verify(context.Background(), tokens)
		require.NoError(st, err)
		require.Equal(st, acc.Username, acc.Username)
		require.Empty(st, acc.PasswordHash)
	})

	t.Run("KO - parse token error", func(st *testing.T) {
		_, err = strategy.Verify(context.Background(), entities.Tokens{})
		require.ErrorContains(st, err, "PASSPORT.UTILS.PARSE_BASIC_CREDENTIALS")
	})

	t.Run("KO - credentials validation error", func(st *testing.T) {
		i := testdata.Fake.IntBetween(0, len(passwords)-1)
		basic := utils.CreateRegionalBasicCredentials(accounts[i].Username + ":")
		tokens := entities.Tokens{Access: utils.SchemeBasic + basic}

		_, err = strategy.Verify(context.Background(), tokens)
		require.ErrorContains(st, err, "PASSPORT.CREDENTIALS")
	})

	t.Run("KO - user not found", func(st *testing.T) {
		basic := utils.CreateRegionalBasicCredentials(
			uuid.NewString() + ":" + testdata.Fake.Internet().Password(),
		)
		tokens := entities.Tokens{Access: utils.SchemeBasic + basic}

		_, err = strategy.Verify(context.Background(), tokens)
		require.ErrorIs(st, err, ErrLogin)
	})

	t.Run("KO - not active error", func(st *testing.T) {
		// setup another batch to test deactivated account
		newaccounts, newpasswords := setup(t)
		require.NoError(st, orm.Create(newaccounts).Error)

		i := testdata.Fake.IntBetween(0, len(passwords)-1)
		user := newaccounts[i].Username
		pass := newpasswords[i]

		err := orm.
			Model(&entities.Account{}).
			Where("username = ?", user).
			Update("deactivated_at", time.Now().Add(-time.Hour).UnixMilli()).
			Error
		require.NoError(st, err)

		basic := utils.CreateRegionalBasicCredentials(user + ":" + pass)
		tokens := entities.Tokens{Access: utils.SchemeBasic + basic}
		_, err = strategy.Verify(context.Background(), tokens)
		require.ErrorIs(st, err, ErrAccountDeactivated)
	})

	t.Run("KO - password not match", func(st *testing.T) {
		i := testdata.Fake.IntBetween(0, len(passwords)-1)
		user := accounts[i].Username
		pass := passwords[i]

		basic := utils.CreateRegionalBasicCredentials(user + ":" + pass + uuid.NewString())
		tokens := entities.Tokens{Access: utils.SchemeBasic + basic}
		_, err := strategy.Verify(context.Background(), tokens)
		require.ErrorIs(st, err, ErrLogin)
	})
}

func TestInternal_Deactivate(t *testing.T) {
	accounts, _ := setup(t)

	strategy, err := NewInternal(internalconf, testify.Logger())
	require.NoError(t, err)

	strategy.Connect(context.Background())
	defer strategy.Disconnect(context.Background())

	orm := strategy.(*internal).orm
	require.NoError(t, orm.Create(accounts).Error)

	t.Run("OK", func(st *testing.T) {
		i := testdata.Fake.IntBetween(0, len(accounts)-1)
		username := accounts[i].Username
		ts := time.Now().UnixMilli()

		err := strategy.Deactivate(context.Background(), username, ts)
		require.NoError(st, err)
	})

	t.Run("KO - user not found", func(st *testing.T) {
		username := uuid.NewString()
		ts := time.Now().UnixMilli()

		err := strategy.Deactivate(context.Background(), username, ts)
		require.ErrorIs(st, err, ErrAccountNotFound)
	})

	t.Run("KO - deactivate timestamp error", func(st *testing.T) {
		i := testdata.Fake.IntBetween(0, len(accounts)-1)
		username := accounts[i].Username
		ts := time.Now().UnixMilli()

		err := strategy.Deactivate(context.Background(), username, ts)
		require.NoError(st, err)

		err = strategy.Deactivate(context.Background(), username, ts-1)
		require.ErrorIs(st, err, ErrDeactivate)
	})
}

func TestInternal_List(t *testing.T) {
	accounts, _ := setup(t)

	strategy, err := NewInternal(internalconf, testify.Logger())
	require.NoError(t, err)

	strategy.Connect(context.Background())
	defer strategy.Disconnect(context.Background())

	orm := strategy.(*internal).orm
	require.NoError(t, orm.Create(accounts).Error)

	t.Run("OK", func(st *testing.T) {
		i := testdata.Fake.IntBetween(0, len(accounts)/2-1)
		j := testdata.Fake.IntBetween(len(accounts)/2, len(accounts)-1)
		usernames := []string{accounts[i].Username, accounts[j].Username}

		acc, err := strategy.List(context.Background(), usernames)
		require.NoError(st, err)

		require.Equal(st, len(usernames), len(acc))
		for i := range acc {
			require.Empty(st, acc[i].PasswordHash)
			require.True(st, slices.Contains(usernames, acc[i].Username))
		}
	})

	t.Run("KO - validation error", func(st *testing.T) {
		i := testdata.Fake.IntBetween(0, len(accounts)/2-1)
		j := testdata.Fake.IntBetween(len(accounts)/2, len(accounts)-1)
		usernames := []string{accounts[i].Username, accounts[j].Username, ""}

		_, err := strategy.List(context.Background(), usernames)
		require.ErrorContains(st, err, fmt.Sprintf("usernames[%d]", len(usernames)-1))
	})
}

func TestInternal_Update(t *testing.T) {
	accounts, _ := setup(t)

	strategy, err := NewInternal(internalconf, testify.Logger())
	require.NoError(t, err)

	strategy.Connect(context.Background())
	defer strategy.Disconnect(context.Background())

	orm := strategy.(*internal).orm
	require.NoError(t, orm.Create(accounts).Error)

	t.Run("OK", func(st *testing.T) {
		i := testdata.Fake.IntBetween(0, len(accounts)-1)
		accounts[i].Name = testdata.Fake.Internet().User()

		require.NoError(st, strategy.Update(context.Background(), accounts[i]))

		accounts[i].Metadata = &safe.Metadata{}
		accounts[i].Metadata.Set("external_id", uuid.NewString())

		require.NoError(st, strategy.Update(context.Background(), accounts[i]))
	})

	t.Run("KO - user not found", func(st *testing.T) {
		err := strategy.Update(context.Background(), entities.Account{Username: uuid.NewString()})
		require.ErrorIs(st, err, ErrAccountNotFound)
	})

	t.Run("KO - not active error", func(st *testing.T) {
		hash, err := password.Hash(uuid.NewString())
		require.NoError(t, err)

		account := entities.Account{
			Username:      uuid.NewString(),
			PasswordHash:  hash,
			Name:          testdata.Fake.Internet().User(),
			CreatedAt:     time.Now().UnixMilli(),
			UpdatedAt:     time.Now().UnixMilli(),
			DeactivatedAt: time.Now().Add(-time.Hour).UnixMilli(),
		}
		require.NoError(t, orm.Create(account).Error)

		err = strategy.Update(context.Background(), account)
		require.ErrorIs(st, err, ErrDeactivate)
	})
}

var internalconf = &config.Internal{Sqlx: sqlxconfig.Config{
	Uri: testdata.SqliteUri,
	Connection: sqlxconfig.Connection{
		MaxLifetime:  sqlxconfig.DefaultConnMaxLifetime,
		MaxIdletime:  sqlxconfig.DefaultConnMaxIdletime,
		MaxIdleCount: sqlxconfig.DefaultConnMaxIdleCount,
		MaxOpenCount: sqlxconfig.DefaultConnMaxOpenCount,
	},
}}
