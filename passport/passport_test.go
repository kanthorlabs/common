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
	sqlx "github.com/kanthorlabs/common/persistence/sqlx/config"
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

		st.Run("KO - Ask init error", func(sst *testing.T) {
			conf := &config.Config{Strategies: make([]config.Strategy, 1)}
			conf.Strategies[0] = ask()
			conf.Strategies[0].Ask.Accounts = append(conf.Strategies[0].Ask.Accounts, conf.Strategies[0].Ask.Accounts[0])

			_, err := New(conf, testify.Logger())
			require.ErrorContains(sst, err, "PASSPORT.STRATEGY.ASK.DUPLICATED_ACCOUNT")
		})

		st.Run("KO - Durability configuration error", func(sst *testing.T) {
			conf := &config.Config{Strategies: make([]config.Strategy, 1)}
			conf.Strategies[0] = durability()
			conf.Strategies[0].Durability = config.Durability{}

			_, err := New(conf, testify.Logger())
			require.ErrorContains(sst, err, "SQLX.CONFIG.")
		})
	})

	t.Run(".Connect/.Readiness/.Liveness/.Disconnect", func(st *testing.T) {
		pp, _ := instance(t)
		require.NoError(st, pp.Connect(context.Background()))
		require.NoError(st, pp.Readiness())
		require.NoError(st, pp.Liveness())
		require.NoError(st, pp.Disconnect(context.Background()))
	})

	t.Run(".Login", func(st *testing.T) {
		t.Run("Ok", func(sst *testing.T) {
			pp, conf := instance(sst)
			require.NoError(sst, pp.Connect(context.Background()))
			defer func() {
				require.NoError(sst, pp.Disconnect(context.Background()))
			}()

			acc := askacc(conf)

			pass, _ := passwords.Load(acc.Username)
			credentials := &entities.Credentials{
				Username: acc.Username,
				Password: pass.(string),
			}
			account, err := pp.Login(context.Background(), askname(conf), credentials)
			require.NoError(sst, err)
			require.Equal(sst, credentials.Username, account.Username)
		})

		t.Run("KO - strategy not found", func(sst *testing.T) {
			pp, _ := instance(sst)
			require.NoError(sst, pp.Connect(context.Background()))
			defer func() {
				require.NoError(sst, pp.Disconnect(context.Background()))
			}()

			_, err := pp.Login(context.Background(), testdata.Fake.Beer().Name(), nil)
			require.ErrorIs(sst, err, ErrStrategyNotFound)
		})
	})

	t.Run(".Logout", func(st *testing.T) {
		t.Run("Ok", func(sst *testing.T) {
			pp, conf := instance(sst)
			require.NoError(sst, pp.Connect(context.Background()))
			defer func() {
				require.NoError(sst, pp.Disconnect(context.Background()))
			}()

			err := pp.Logout(context.Background(), askname(conf), nil)
			require.NoError(sst, err)
		})

		t.Run("KO - strategy not found", func(sst *testing.T) {
			pp, _ := instance(sst)
			require.NoError(sst, pp.Connect(context.Background()))
			defer func() {
				require.NoError(sst, pp.Disconnect(context.Background()))
			}()

			err := pp.Logout(context.Background(), testdata.Fake.Beer().Name(), nil)
			require.ErrorIs(sst, err, ErrStrategyNotFound)
		})
	})

	t.Run(".Verify", func(st *testing.T) {
		t.Run("Ok", func(sst *testing.T) {
			pp, conf := instance(sst)
			require.NoError(sst, pp.Connect(context.Background()))
			defer func() {
				require.NoError(sst, pp.Disconnect(context.Background()))
			}()

			acc := askacc(conf)
			pass, _ := passwords.Load(acc.Username)
			credentials := &entities.Credentials{
				Username: acc.Username,
				Password: pass.(string),
			}
			account, err := pp.Verify(context.Background(), askname(conf), credentials)
			require.NoError(sst, err)
			require.Equal(sst, credentials.Username, account.Username)
		})

		t.Run("KO - strategy not found", func(sst *testing.T) {
			pp, _ := instance(sst)
			require.NoError(sst, pp.Connect(context.Background()))
			defer func() {
				require.NoError(sst, pp.Disconnect(context.Background()))
			}()

			_, err := pp.Verify(context.Background(), testdata.Fake.Beer().Name(), nil)
			require.ErrorIs(sst, err, ErrStrategyNotFound)
		})
	})

	t.Run(".Register", func(st *testing.T) {
		t.Run("OK", func(sst *testing.T) {
			pp, conf := instance(sst)
			require.NoError(sst, pp.Connect(context.Background()))
			defer func() {
				require.NoError(sst, pp.Disconnect(context.Background()))
			}()

			err := pp.Register(context.Background(), askname(conf), nil)
			require.ErrorContains(sst, err, "PASSPORT.ASK.REGISTER.UNIMPLEMENT")
		})

		t.Run("KO - strategy not found", func(sst *testing.T) {
			pp, _ := instance(sst)
			require.NoError(sst, pp.Connect(context.Background()))
			defer func() {
				require.NoError(sst, pp.Disconnect(context.Background()))
			}()

			err := pp.Register(context.Background(), testdata.Fake.Beer().Name(), nil)
			require.ErrorIs(sst, err, ErrStrategyNotFound)
		})
	})
}

func instance(t *testing.T) (Passport, *config.Config) {
	conf := &config.Config{Strategies: make([]config.Strategy, 0)}
	conf.Strategies = append(conf.Strategies, durability())
	conf.Strategies = append(conf.Strategies, ask())

	pp, err := New(conf, testify.Logger())
	require.NoError(t, err)

	return pp, conf
}

func ask() config.Strategy {
	pass := testdata.Fake.Internet().Password()
	hash, _ := password.HashString(pass)
	account := entities.Account{
		Username:     uuid.NewString(),
		PasswordHash: hash,
		Name:         testdata.Fake.Internet().User(),
		CreatedAt:    time.Now().UnixMilli(),
		UpdatedAt:    time.Now().UnixMilli(),
	}

	passwords.Store(account.Username, pass)

	return config.Strategy{
		Engine: config.EngineAsk,
		Name:   uuid.NewString(),
		Ask: config.Ask{
			Accounts: []entities.Account{account},
		},
	}
}

func askacc(conf *config.Config) entities.Account {
	for i := range conf.Strategies {
		if conf.Strategies[i].Engine == config.EngineAsk {
			j := testdata.Fake.IntBetween(0, len(conf.Strategies[i].Ask.Accounts)-1)
			return conf.Strategies[i].Ask.Accounts[j]
		}
	}
	panic("no ask strategy was configured")
}

func askname(conf *config.Config) string {
	for i := range conf.Strategies {
		if conf.Strategies[i].Engine == config.EngineAsk {
			return conf.Strategies[i].Name
		}
	}
	panic("no ask strategy was configured")
}

func durability() config.Strategy {
	return config.Strategy{
		Engine: config.EngineDurability,
		Name:   uuid.NewString(),
		Durability: config.Durability{
			Sqlx: sqlx.Config{
				Uri: testdata.SqliteUri,
				Connection: sqlx.Connection{
					MaxLifetime:  sqlx.DefaultConnMaxLifetime,
					MaxIdletime:  sqlx.DefaultConnMaxIdletime,
					MaxIdleCount: sqlx.DefaultConnMaxIdleCount,
					MaxOpenCount: sqlx.DefaultConnMaxOpenCount,
				},
			},
		},
	}
}
