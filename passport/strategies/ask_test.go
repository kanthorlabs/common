package strategies

import (
	"context"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/passport/config"
	"github.com/kanthorlabs/common/passport/entities"
	"github.com/kanthorlabs/common/passport/utils"
	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/require"
)

func TestAsk_New(t *testing.T) {
	accounts, _ := setup(t)

	t.Run("OK", func(st *testing.T) {
		conf := &config.Ask{Accounts: accounts}
		_, err := NewAsk(conf, testify.Logger())
		require.NoError(st, err)
	})

	t.Run("KO - configuration error", func(st *testing.T) {
		conf := &config.Ask{}
		_, err := NewAsk(conf, testify.Logger())
		require.ErrorContains(st, err, "PASSPORT.STRATEGY.ASK.CONFIG")
	})

	t.Run("KO - duplicated account error", func(st *testing.T) {
		conf := &config.Ask{
			Accounts: []entities.Account{
				{
					Username:     uuid.NewString(),
					PasswordHash: testdata.Fake.Internet().Password(),
					Name:         testdata.Fake.Internet().User(),
					CreatedAt:    time.Now().UnixMilli(),
					UpdatedAt:    time.Now().UnixMilli(),
				},
			},
		}
		conf.Accounts = append(conf.Accounts, conf.Accounts[0])
		_, err := NewAsk(conf, testify.Logger())
		require.ErrorContains(st, err, "PASSPORT.STRATEGY.ASK.DUPLICATED_ACCOUNT")
	})
}

func TestAsk_Connect(t *testing.T) {
	accounts, _ := setup(t)

	t.Run("OK", func(st *testing.T) {
		conf := &config.Ask{Accounts: accounts}
		strategy, err := NewAsk(conf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, strategy.Connect(context.Background()))
	})

	t.Run("KO - already connected error", func(st *testing.T) {
		conf := &config.Ask{Accounts: accounts}
		strategy, err := NewAsk(conf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, strategy.Connect(context.Background()))
		require.ErrorIs(st, strategy.Connect(context.Background()), ErrAlreadyConnected)
	})
}

func TestAsk_Readiness(t *testing.T) {
	accounts, _ := setup(t)

	t.Run("OK", func(st *testing.T) {
		conf := &config.Ask{Accounts: accounts}
		strategy, err := NewAsk(conf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, strategy.Connect(context.Background()))
		require.NoError(st, strategy.Readiness())
	})

	t.Run("OK - disconnected", func(st *testing.T) {
		conf := &config.Ask{Accounts: accounts}
		strategy, err := NewAsk(conf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, strategy.Connect(context.Background()))
		require.NoError(st, strategy.Disconnect(context.Background()))
		require.NoError(st, strategy.Readiness())
	})

	t.Run("KO - not connected error", func(st *testing.T) {
		conf := &config.Ask{Accounts: accounts}
		strategy, err := NewAsk(conf, testify.Logger())
		require.NoError(st, err)

		require.ErrorIs(st, strategy.Readiness(), ErrNotConnected)
	})
}

func TestAsk_Liveness(t *testing.T) {
	accounts, _ := setup(t)

	t.Run("OK", func(st *testing.T) {
		conf := &config.Ask{Accounts: accounts}
		strategy, err := NewAsk(conf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, strategy.Connect(context.Background()))
		require.NoError(st, strategy.Liveness())
	})

	t.Run("OK - disconnected", func(st *testing.T) {
		conf := &config.Ask{Accounts: accounts}
		strategy, err := NewAsk(conf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, strategy.Connect(context.Background()))
		require.NoError(st, strategy.Disconnect(context.Background()))
		require.NoError(st, strategy.Liveness())
	})

	t.Run("KO - not connected error", func(st *testing.T) {
		conf := &config.Ask{Accounts: accounts}
		strategy, err := NewAsk(conf, testify.Logger())
		require.NoError(st, err)

		require.ErrorIs(st, strategy.Liveness(), ErrNotConnected)
	})
}

func TestAsk_Disconnect(t *testing.T) {
	accounts, _ := setup(t)

	t.Run("OK", func(st *testing.T) {
		conf := &config.Ask{Accounts: accounts}
		strategy, err := NewAsk(conf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, strategy.Connect(context.Background()))
		require.NoError(st, strategy.Disconnect(context.Background()))
	})

	t.Run("KO - not connected error", func(st *testing.T) {
		conf := &config.Ask{Accounts: accounts}
		strategy, err := NewAsk(conf, testify.Logger())
		require.NoError(st, err)

		require.ErrorIs(st, strategy.Disconnect(context.Background()), ErrNotConnected)
	})
}

func TestAsk_Register(t *testing.T) {
	accounts, _ := setup(t)

	conf := &config.Ask{Accounts: accounts}
	strategy, err := NewAsk(conf, testify.Logger())
	require.NoError(t, err)

	t.Run("KO - unimplement error", func(st *testing.T) {
		err = strategy.Register(context.Background(), accounts[0])
		require.ErrorContains(st, err, "UNIMPLEMENT.ERROR")
	})
}

func TestAsk_Login(t *testing.T) {
	accounts, passwords := setup(t)

	conf := &config.Ask{Accounts: accounts}
	strategy, err := NewAsk(conf, testify.Logger())
	require.NoError(t, err)

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

func TestAsk_Logout(t *testing.T) {
	accounts, _ := setup(t)

	conf := &config.Ask{Accounts: accounts}
	strategy, err := NewAsk(conf, testify.Logger())
	require.NoError(t, err)

	t.Run("KO - unimplement error", func(st *testing.T) {
		err := strategy.Logout(context.Background(), entities.Tokens{})
		require.ErrorContains(st, err, "UNIMPLEMENT.ERROR")
	})
}

func TestAsk_Verify(t *testing.T) {
	accounts, passwords := setup(t)

	conf := &config.Ask{Accounts: accounts}
	strategy, err := NewAsk(conf, testify.Logger())
	require.NoError(t, err)

	t.Run("OK", func(st *testing.T) {
		i := testdata.Fake.IntBetween(0, len(passwords)-1)
		user := accounts[i].Username
		pass := passwords[i]

		tokens := entities.Tokens{
			Access: utils.SchemeBasic + utils.CreateRegionalBasicCredentials(user+":"+pass),
		}
		acc, err := strategy.Verify(context.Background(), tokens)
		require.NoError(st, err)
		require.Equal(st, user, acc.Username)
		require.Empty(st, acc.PasswordHash)
	})

	t.Run("KO - parse token error", func(st *testing.T) {
		_, err = strategy.Verify(context.Background(), entities.Tokens{})
		require.ErrorContains(st, err, "PASSPORT.UTILS.PARSE_BASIC_CREDENTIALS")
	})

	t.Run("KO - credentials validation error", func(st *testing.T) {
		i := testdata.Fake.IntBetween(0, len(passwords)-1)
		tokens := entities.Tokens{
			Access: utils.SchemeBasic + utils.CreateRegionalBasicCredentials(":"+passwords[i]),
		}

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

	t.Run("KO - password not match", func(st *testing.T) {
		i := testdata.Fake.IntBetween(0, len(passwords)-1)
		user := accounts[i].Username
		pass := passwords[i]

		basic := utils.CreateRegionalBasicCredentials(user + ":" + pass + uuid.NewString())
		tokens := entities.Tokens{
			Access: utils.SchemeBasic + basic,
		}
		_, err = strategy.Verify(context.Background(), tokens)
		require.ErrorIs(st, err, ErrLogin)
	})
}

func TestAsk_Deactivate(t *testing.T) {
	accounts, _ := setup(t)

	conf := &config.Ask{Accounts: accounts}
	strategy, err := NewAsk(conf, testify.Logger())
	require.NoError(t, err)

	t.Run("KO - unimplement error", func(st *testing.T) {
		err = strategy.Deactivate(context.Background(), accounts[0].Username, time.Now().UnixMilli())
		require.ErrorContains(st, err, "UNIMPLEMENT.ERROR")
	})
}

func TestAsk_List(t *testing.T) {
	accounts, _ := setup(t)

	conf := &config.Ask{Accounts: accounts}
	strategy, err := NewAsk(conf, testify.Logger())
	require.NoError(t, err)

	t.Run("OK", func(st *testing.T) {
		usernames := []string{accounts[0].Username, accounts[1].Username}
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

func TestAsk_Update(t *testing.T) {
	accounts, _ := setup(t)

	conf := &config.Ask{Accounts: accounts}
	strategy, err := NewAsk(conf, testify.Logger())
	require.NoError(t, err)

	t.Run("KO - unimplement error", func(st *testing.T) {
		err := strategy.Update(context.Background(), entities.Account{})
		require.ErrorContains(st, err, "UNIMPLEMENT.ERROR")
	})
}
