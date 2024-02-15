package passport

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/cipher/password"
	"github.com/kanthorlabs/common/passport/config"
	"github.com/kanthorlabs/common/passport/entities"
	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/require"
)

var passwords = sync.Map{}

func TestPassport(t *testing.T) {
	t.Run("New", func(st *testing.T) {
		st.Run("KO - configuration error", func(sst *testing.T) {
			_, err := New(&config.Config{}, testify.Logger())
			require.ErrorContains(sst, err, "PASSPORT.CONFIG")
		})

		st.Run("KO - duplicated name", func(sst *testing.T) {
			conf := &config.Config{Strategies: make([]config.Strategy, 2)}
			conf.Strategies[0] = ask()
			conf.Strategies[1] = ask()
			conf.Strategies[1].Name = conf.Strategies[0].Name
			_, err := New(conf, testify.Logger())
			require.ErrorIs(sst, err, ErrStrategyDuplicated)
		})

		st.Run("KO - Ask configuration error", func(sst *testing.T) {
			conf := &config.Config{Strategies: make([]config.Strategy, 1)}
			conf.Strategies[0] = ask()
			conf.Strategies[0].Ask.Accounts = make([]entities.Account, 0)

			_, err := New(conf, testify.Logger())
			require.ErrorContains(sst, err, "PASSPORT.STRATEGY.ASK.CONFIG")
		})
	})

	t.Run(".Connect/.Readiness/.Liveness/.Disconnect", func(st *testing.T) {
		pp, _ := instance(t)
		require.Nil(st, pp.Connect(context.Background()))
		require.Nil(st, pp.Readiness())
		require.Nil(st, pp.Liveness())
		require.Nil(st, pp.Disconnect(context.Background()))
	})

	t.Run(".Login", func(st *testing.T) {
		t.Run("Ok", func(sst *testing.T) {
			pp, conf := instance(sst)
			require.Nil(sst, pp.Connect(context.Background()))
			defer func() {
				require.Nil(sst, pp.Disconnect(context.Background()))
			}()

			pass, _ := passwords.Load(conf.Strategies[0].Ask.Accounts[0].Username)
			credentials := &entities.Credentials{
				Username: conf.Strategies[0].Ask.Accounts[0].Username,
				Password: pass.(string),
			}
			acc, err := pp.Login(context.Background(), conf.Strategies[0].Name, credentials)
			require.Nil(sst, err)
			require.Equal(sst, credentials.Username, acc.Username)
		})

		t.Run("KO - strategy not found", func(sst *testing.T) {
			pp, _ := instance(sst)
			require.Nil(sst, pp.Connect(context.Background()))
			defer func() {
				require.Nil(sst, pp.Disconnect(context.Background()))
			}()

			_, err := pp.Login(context.Background(), testdata.Fake.Beer().Name(), nil)
			require.ErrorIs(sst, err, ErrStrategyNotFound)
		})
	})

	t.Run(".Logout", func(st *testing.T) {
		t.Run("Ok", func(sst *testing.T) {
			pp, conf := instance(sst)
			require.Nil(sst, pp.Connect(context.Background()))
			defer func() {
				require.Nil(sst, pp.Disconnect(context.Background()))
			}()

			err := pp.Logout(context.Background(), conf.Strategies[0].Name, nil)
			require.Nil(sst, err)
		})

		t.Run("KO - strategy not found", func(sst *testing.T) {
			pp, _ := instance(sst)
			require.Nil(sst, pp.Connect(context.Background()))
			defer func() {
				require.Nil(sst, pp.Disconnect(context.Background()))
			}()

			err := pp.Logout(context.Background(), testdata.Fake.Beer().Name(), nil)
			require.ErrorIs(sst, err, ErrStrategyNotFound)
		})
	})

	t.Run(".Verify", func(st *testing.T) {
		t.Run("Ok", func(sst *testing.T) {
			pp, conf := instance(sst)
			require.Nil(sst, pp.Connect(context.Background()))
			defer func() {
				require.Nil(sst, pp.Disconnect(context.Background()))
			}()

			pass, _ := passwords.Load(conf.Strategies[0].Ask.Accounts[0].Username)
			credentials := &entities.Credentials{
				Username: conf.Strategies[0].Ask.Accounts[0].Username,
				Password: pass.(string),
			}
			acc, err := pp.Verify(context.Background(), conf.Strategies[0].Name, credentials)
			require.Nil(sst, err)
			require.Equal(sst, credentials.Username, acc.Username)
		})

		t.Run("KO - strategy not found", func(sst *testing.T) {
			pp, _ := instance(sst)
			require.Nil(sst, pp.Connect(context.Background()))
			defer func() {
				require.Nil(sst, pp.Disconnect(context.Background()))
			}()

			_, err := pp.Verify(context.Background(), testdata.Fake.Beer().Name(), nil)
			require.ErrorIs(sst, err, ErrStrategyNotFound)
		})
	})

	t.Run(".Register", func(st *testing.T) {
		t.Run("OK", func(sst *testing.T) {
			pp, conf := instance(sst)
			require.Nil(sst, pp.Connect(context.Background()))
			defer func() {
				require.Nil(sst, pp.Disconnect(context.Background()))
			}()

			err := pp.Register(context.Background(), conf.Strategies[0].Name, nil)
			require.ErrorContains(sst, err, "PASSPORT.ASK.REGISTER.UNIMPLEMENT")
		})

		t.Run("KO - strategy not found", func(sst *testing.T) {
			pp, _ := instance(sst)
			require.Nil(sst, pp.Connect(context.Background()))
			defer func() {
				require.Nil(sst, pp.Disconnect(context.Background()))
			}()

			err := pp.Register(context.Background(), testdata.Fake.Beer().Name(), nil)
			require.ErrorIs(sst, err, ErrStrategyNotFound)
		})
	})
}

func instance(t *testing.T) (Passport, *config.Config) {
	conf := &config.Config{Strategies: make([]config.Strategy, 0)}
	conf.Strategies = append(conf.Strategies, ask())

	pp, err := New(conf, testify.Logger())
	require.Nil(t, err)

	return pp, conf
}

func ask() config.Strategy {
	pass := testdata.Fake.Internet().Password()
	hash, _ := password.HashString(pass)
	account := entities.Account{
		Username:     uuid.NewString(),
		PasswordHash: hash,
		Tenant:       entities.TenantSuper,
		Name:         testdata.Fake.Internet().User(),
		CreatedAt:    time.Now().UnixMilli(),
		UpdatedAt:    time.Now().UnixMilli(),
	}

	passwords.Store(account.Username, pass)

	return config.Strategy{
		Engine: config.EngineAsk,
		Name:   testdata.Fake.Beer().Name(),
		Ask: config.Ask{
			Accounts: []entities.Account{account},
		},
	}
}
