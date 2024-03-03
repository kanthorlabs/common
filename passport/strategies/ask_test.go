package strategies

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/passport/config"
	"github.com/kanthorlabs/common/passport/entities"
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
		c, err := NewAsk(conf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, c.Connect(context.Background()))
	})

	t.Run("KO - already connected", func(st *testing.T) {
		conf := &config.Ask{Accounts: accounts}
		c, err := NewAsk(conf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, c.Connect(context.Background()))
		require.ErrorIs(st, c.Connect(context.Background()), ErrAlreadyConnected)
	})
}

func TestAsk_Readiness(t *testing.T) {
	accounts, _ := setup(t)

	t.Run("OK", func(st *testing.T) {
		conf := &config.Ask{Accounts: accounts}
		c, err := NewAsk(conf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, c.Connect(context.Background()))
		require.NoError(st, c.Readiness())
	})

	t.Run("OK - disconnected", func(st *testing.T) {
		conf := &config.Ask{Accounts: accounts}
		c, err := NewAsk(conf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, c.Connect(context.Background()))
		require.NoError(st, c.Disconnect(context.Background()))
		require.NoError(st, c.Readiness())
	})

	t.Run("KO - not connected", func(st *testing.T) {
		conf := &config.Ask{Accounts: accounts}
		c, err := NewAsk(conf, testify.Logger())
		require.NoError(st, err)

		require.ErrorIs(st, c.Readiness(), ErrNotConnected)
	})
}

func TestAsk_Liveness(t *testing.T) {
	accounts, _ := setup(t)

	t.Run("OK", func(st *testing.T) {
		conf := &config.Ask{Accounts: accounts}
		c, err := NewAsk(conf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, c.Connect(context.Background()))
		require.NoError(st, c.Liveness())
	})

	t.Run("OK - disconnected", func(st *testing.T) {
		conf := &config.Ask{Accounts: accounts}
		c, err := NewAsk(conf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, c.Connect(context.Background()))
		require.NoError(st, c.Disconnect(context.Background()))
		require.NoError(st, c.Liveness())
	})

	t.Run("KO - not connected", func(st *testing.T) {
		conf := &config.Ask{Accounts: accounts}
		c, err := NewAsk(conf, testify.Logger())
		require.NoError(st, err)

		require.ErrorIs(st, c.Liveness(), ErrNotConnected)
	})
}

func TestAsk_Login(t *testing.T) {
	accounts, passwords := setup(t)

	t.Run("OK", func(st *testing.T) {
		conf := &config.Ask{Accounts: accounts}
		c, err := NewAsk(conf, testify.Logger())
		require.NoError(st, err)

		i := testdata.Fake.IntBetween(0, len(passwords)-1)
		credentials := &entities.Credentials{
			Username: accounts[i].Username,
			Password: passwords[i],
		}
		acc, err := c.Login(context.Background(), credentials)
		require.NoError(st, err)
		require.Equal(st, credentials.Username, acc.Username)
		require.Empty(st, acc.PasswordHash)
	})

	t.Run("KO - credentials error", func(st *testing.T) {
		conf := &config.Ask{Accounts: accounts}
		c, err := NewAsk(conf, testify.Logger())
		require.NoError(st, err)

		_, err = c.Login(context.Background(), nil)
		require.ErrorContains(st, err, "PASSPORT.CREDENTIALS")

		_, err = c.Login(context.Background(), &entities.Credentials{})
		require.ErrorContains(st, err, "PASSPORT.CREDENTIALS")
	})

	t.Run("KO - user not found", func(st *testing.T) {
		conf := &config.Ask{Accounts: accounts}
		c, err := NewAsk(conf, testify.Logger())
		require.NoError(st, err)

		credentials := &entities.Credentials{
			Username: uuid.NewString(),
			Password: testdata.Fake.Internet().Password(),
		}
		_, err = c.Login(context.Background(), credentials)
		require.ErrorIs(st, err, ErrLogin)
	})

	t.Run("KO - password not match", func(st *testing.T) {
		conf := &config.Ask{Accounts: accounts}
		c, err := NewAsk(conf, testify.Logger())
		require.NoError(st, err)

		i := testdata.Fake.IntBetween(0, len(passwords)-1)
		credentials := &entities.Credentials{
			Username: accounts[i].Username,
			Password: testdata.Fake.Internet().Password(),
		}
		_, err = c.Login(context.Background(), credentials)
		require.ErrorIs(st, err, ErrLogin)
	})
}

func TestAsk_Verify(t *testing.T) {
	accounts, passwords := setup(t)

	t.Run("OK", func(st *testing.T) {
		conf := &config.Ask{Accounts: accounts}
		c, err := NewAsk(conf, testify.Logger())
		require.NoError(st, err)

		i := testdata.Fake.IntBetween(0, len(passwords)-1)
		credentials := &entities.Credentials{
			Username: accounts[i].Username,
			Password: passwords[i],
		}
		acc, err := c.Verify(context.Background(), credentials)
		require.NoError(st, err)
		require.Equal(st, credentials.Username, acc.Username)
		require.Empty(st, acc.PasswordHash)
	})

	t.Run("KO - credentials error", func(st *testing.T) {
		conf := &config.Ask{Accounts: accounts}
		c, err := NewAsk(conf, testify.Logger())
		require.NoError(st, err)

		_, err = c.Verify(context.Background(), nil)
		require.ErrorContains(st, err, "PASSPORT.CREDENTIALS")

		_, err = c.Verify(context.Background(), &entities.Credentials{})
		require.ErrorContains(st, err, "PASSPORT.CREDENTIALS")
	})

	t.Run("KO - user not found", func(st *testing.T) {
		conf := &config.Ask{Accounts: accounts}
		c, err := NewAsk(conf, testify.Logger())
		require.NoError(st, err)

		credentials := &entities.Credentials{
			Username: uuid.NewString(),
			Password: testdata.Fake.Internet().Password(),
		}
		_, err = c.Verify(context.Background(), credentials)
		require.ErrorIs(st, err, ErrLogin)
	})

	t.Run("KO - password not match", func(st *testing.T) {
		conf := &config.Ask{Accounts: accounts}
		c, err := NewAsk(conf, testify.Logger())
		require.NoError(st, err)

		i := testdata.Fake.IntBetween(0, len(passwords)-1)
		credentials := &entities.Credentials{
			Username: accounts[i].Username,
			Password: testdata.Fake.Internet().Password(),
		}
		_, err = c.Verify(context.Background(), credentials)
		require.ErrorIs(st, err, ErrLogin)
	})
}

func TestAsk_Logout(t *testing.T) {
	accounts, _ := setup(t)

	t.Run("OK", func(st *testing.T) {
		conf := &config.Ask{Accounts: accounts}
		c, err := NewAsk(conf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, c.Logout(context.Background(), &entities.Credentials{}))
	})
}

func TestAsk_Register(t *testing.T) {
	accounts, _ := setup(t)

	t.Run("KO - unimplement error", func(st *testing.T) {
		conf := &config.Ask{Accounts: accounts}
		c, err := NewAsk(conf, testify.Logger())
		require.NoError(st, err)

		err = c.Register(context.Background(), &accounts[0])
		require.ErrorContains(st, err, "PASSPORT.ASK.REGISTER.UNIMPLEMENT")
	})
}

func TestAsk_Deactivate(t *testing.T) {
	accounts, _ := setup(t)

	t.Run("KO - unimplement error", func(st *testing.T) {
		conf := &config.Ask{Accounts: accounts}
		c, err := NewAsk(conf, testify.Logger())
		require.NoError(st, err)

		err = c.Deactivate(context.Background(), accounts[0].Username, time.Now().UnixMilli())
		require.ErrorContains(st, err, "PASSPORT.ASK.DEACTIVATE.UNIMPLEMENT")
	})
}
