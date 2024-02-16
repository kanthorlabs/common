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

func TestAsk(t *testing.T) {
	accounts, passwords := setup(t)

	t.Run("New", func(st *testing.T) {
		st.Run("KO - configuration error", func(sst *testing.T) {
			_, err := NewAsk(&config.Ask{}, testify.Logger())
			require.ErrorContains(sst, err, "PASSPORT.STRATEGY.ASK.CONFIG")
		})

		st.Run("KO - duplicated account error", func(sst *testing.T) {
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
			require.ErrorContains(sst, err, "PASSPORT.STRATEGY.ASK.DUPLICATED_ACCOUNT")
		})
	})

	t.Run(".Connect/.Readiness/.Liveness/.Logout/.Disconnect", func(st *testing.T) {
		conf := &config.Ask{Accounts: accounts}
		strategy, err := NewAsk(conf, testify.Logger())
		require.Nil(st, err)

		require.Nil(st, strategy.Connect(context.Background()))
		require.Nil(st, strategy.Readiness())
		require.Nil(st, strategy.Liveness())
		require.Nil(st, strategy.Logout(context.Background(), nil))
		require.Nil(st, strategy.Disconnect(context.Background()))
	})

	t.Run(".Login", func(st *testing.T) {
		conf := &config.Ask{Accounts: accounts}
		strategy, err := NewAsk(conf, testify.Logger())
		require.Nil(st, err)

		st.Run("OK", func(sst *testing.T) {
			i := testdata.Fake.IntBetween(0, len(passwords)-1)
			credentials := &entities.Credentials{
				Username: accounts[i].Username,
				Password: passwords[i],
			}
			acc, err := strategy.Login(context.Background(), credentials)
			require.Nil(sst, err)
			require.Equal(sst, credentials.Username, acc.Username)
			require.Empty(sst, acc.PasswordHash)
		})

		st.Run("KO - credentials error", func(sst *testing.T) {
			_, err := strategy.Login(context.Background(), nil)
			require.ErrorContains(sst, err, "PASSPORT.CREDENTIALS")

			_, err = strategy.Login(context.Background(), &entities.Credentials{})
			require.ErrorContains(sst, err, "PASSPORT.CREDENTIALS")
		})

		st.Run("KO - user not found", func(sst *testing.T) {
			credentials := &entities.Credentials{
				Username: uuid.NewString(),
				Password: testdata.Fake.Internet().Password(),
			}
			_, err := strategy.Login(context.Background(), credentials)
			require.ErrorIs(sst, err, ErrLogin)
		})

		st.Run("KO - password not match", func(sst *testing.T) {
			i := testdata.Fake.IntBetween(0, len(passwords)-1)
			credentials := &entities.Credentials{
				Username: accounts[i].Username,
				Password: testdata.Fake.Internet().Password(),
			}
			_, err := strategy.Login(context.Background(), credentials)
			require.ErrorIs(sst, err, ErrLogin)
		})
	})

	t.Run(".Verify", func(st *testing.T) {
		conf := &config.Ask{Accounts: accounts}
		strategy, err := NewAsk(conf, testify.Logger())
		require.Nil(st, err)

		st.Run("OK", func(sst *testing.T) {
			i := testdata.Fake.IntBetween(0, len(passwords)-1)
			credentials := &entities.Credentials{
				Username: accounts[i].Username,
				Password: passwords[i],
			}
			acc, err := strategy.Verify(context.Background(), credentials)
			require.Nil(sst, err)
			require.Equal(sst, credentials.Username, acc.Username)
			require.Empty(sst, acc.PasswordHash)
		})
	})

	t.Run(".Register", func(st *testing.T) {
		conf := &config.Ask{Accounts: accounts}
		strategy, err := NewAsk(conf, testify.Logger())
		require.Nil(st, err)

		st.Run("KO - unimplement error", func(sst *testing.T) {
			i := testdata.Fake.IntBetween(0, len(passwords)-1)
			err := strategy.Register(context.Background(), &accounts[i])
			require.ErrorContains(sst, err, "PASSPORT.ASK.REGISTER.UNIMPLEMENT")
		})
	})
}
