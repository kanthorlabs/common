package passport

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/cipher/password"
	"github.com/kanthorlabs/common/passport/config"
	"github.com/kanthorlabs/common/passport/entities"
	"github.com/kanthorlabs/common/passport/utils"
	sqlxconfig "github.com/kanthorlabs/common/persistence/sqlx/config"
	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/require"
)

var passwords = sync.Map{}

func TestPassport_New(t *testing.T) {
	t.Run("KO - configuration error", func(st *testing.T) {
		_, err := New(&config.Config{}, testify.Logger())
		require.ErrorContains(st, err, "PASSPORT.CONFIG")
	})

	t.Run("KO - duplicated name", func(st *testing.T) {
		conf := &config.Config{Strategies: make([]config.Strategy, 2)}
		conf.Strategies[0] = ask()
		conf.Strategies[1] = ask()
		conf.Strategies[1].Name = conf.Strategies[0].Name
		_, err := New(conf, testify.Logger())
		require.ErrorIs(st, err, ErrStrategyDuplicated)
	})

	t.Run("KO - Ask configuration error", func(st *testing.T) {
		conf := &config.Config{Strategies: make([]config.Strategy, 1)}
		conf.Strategies[0] = ask()
		conf.Strategies[0].Ask.Accounts = make([]entities.Account, 0)

		_, err := New(conf, testify.Logger())
		require.ErrorContains(st, err, "PASSPORT.STRATEGY.ASK.CONFIG")
	})

	t.Run("KO - Ask init error", func(st *testing.T) {
		conf := &config.Config{Strategies: make([]config.Strategy, 1)}
		conf.Strategies[0] = ask()
		conf.Strategies[0].Ask.Accounts = append(conf.Strategies[0].Ask.Accounts, conf.Strategies[0].Ask.Accounts[0])

		_, err := New(conf, testify.Logger())
		require.ErrorContains(st, err, "PASSPORT.STRATEGY.ASK.DUPLICATED_ACCOUNT")
	})

	t.Run("KO - Internal configuration error", func(st *testing.T) {
		conf := &config.Config{Strategies: make([]config.Strategy, 1)}
		conf.Strategies[0] = internal()
		conf.Strategies[0].Internal = config.Internal{}

		_, err := New(conf, testify.Logger())
		require.ErrorContains(st, err, "SQLX.CONFIG.")
	})

	t.Run("KO - External configuration error", func(st *testing.T) {
		server := httptestserver()
		defer server.Close()

		conf := &config.Config{Strategies: make([]config.Strategy, 1)}
		conf.Strategies[0] = external(server.URL)
		conf.Strategies[0].External = config.External{}

		_, err := New(conf, testify.Logger())
		require.ErrorContains(st, err, "PASSPORT.CONFIG.EXTERNAL.")
	})
}

func TestPassport_Connect(t *testing.T) {
	server := httptestserver()
	defer server.Close()

	t.Run("OK", func(st *testing.T) {
		pp, _ := instance(t, server)

		require.NoError(st, pp.Connect(context.Background()))
	})

	t.Run("KO - already connected error", func(st *testing.T) {
		pp, _ := instance(t, server)

		require.NoError(st, pp.Connect(context.Background()))
		require.ErrorIs(st, pp.Connect(context.Background()), ErrAlreadyConnected)
	})
}

func TestPassport_Readiness(t *testing.T) {
	server := httptestserver()
	defer server.Close()

	t.Run("OK", func(st *testing.T) {
		pp, _ := instance(t, server)

		require.NoError(st, pp.Connect(context.Background()))
		require.NoError(st, pp.Readiness())
	})

	t.Run("OK - disconnected", func(st *testing.T) {
		pp, _ := instance(t, server)

		require.NoError(st, pp.Connect(context.Background()))
		require.NoError(st, pp.Disconnect(context.Background()))

		require.NoError(st, pp.Readiness())
	})

	t.Run("KO - not connected error", func(st *testing.T) {
		pp, _ := instance(t, server)

		require.ErrorIs(st, pp.Readiness(), ErrNotConnected)
	})
}

func TestPassport_Liveness(t *testing.T) {
	server := httptestserver()
	defer server.Close()

	t.Run("OK", func(st *testing.T) {
		pp, _ := instance(t, server)

		require.NoError(st, pp.Connect(context.Background()))
		require.NoError(st, pp.Liveness())
	})

	t.Run("OK - disconnected", func(st *testing.T) {
		pp, _ := instance(t, server)

		require.NoError(st, pp.Connect(context.Background()))
		require.NoError(st, pp.Disconnect(context.Background()))

		require.NoError(st, pp.Liveness())
	})
	t.Run("KO - not connected error", func(st *testing.T) {
		pp, _ := instance(t, server)

		require.ErrorIs(st, pp.Liveness(), ErrNotConnected)
	})
}

func TestPassport_Disconnect(t *testing.T) {
	server := httptestserver()
	defer server.Close()

	t.Run("OK", func(st *testing.T) {
		pp, _ := instance(t, server)

		require.NoError(st, pp.Connect(context.Background()))
		require.NoError(st, pp.Disconnect(context.Background()))
	})

	t.Run("KO - not connected error", func(st *testing.T) {
		pp, _ := instance(t, server)

		require.ErrorIs(st, pp.Disconnect(context.Background()), ErrNotConnected)
	})
}

func TestPassport_Strategy(t *testing.T) {
	server := httptestserver()
	defer server.Close()

	t.Run("OK", func(st *testing.T) {
		pp, conf := instance(t, server)

		require.NoError(st, pp.Connect(context.Background()))
		defer func() {
			require.NoError(st, pp.Disconnect(context.Background()))
		}()

		strategy, err := pp.Strategy(askname(conf))
		require.NoError(st, err)
		account := askacc(conf)

		pass, _ := passwords.Load(account.Username)
		tokens := entities.Tokens{
			Access: utils.SchemeBasic + utils.CreateRegionalBasicCredentials(account.Username+":"+pass.(string)),
		}
		acc, err := strategy.Verify(context.Background(), tokens)
		require.NoError(st, err)
		require.Equal(st, account.Username, acc.Username)
		require.Empty(st, acc.PasswordHash)
	})

	t.Run("KO - strategy not found", func(st *testing.T) {
		pp, _ := instance(t, server)

		require.NoError(st, pp.Connect(context.Background()))
		defer func() {
			require.NoError(st, pp.Disconnect(context.Background()))
		}()

		_, err := pp.Strategy(testdata.Fake.Beer().Name())
		require.ErrorIs(st, err, ErrStrategyNotFound)
	})
}

func instance(t *testing.T, server *httptest.Server) (Passport, *config.Config) {
	conf := &config.Config{Strategies: make([]config.Strategy, 0)}
	conf.Strategies = append(conf.Strategies, internal())
	conf.Strategies = append(conf.Strategies, ask())
	conf.Strategies = append(conf.Strategies, external(server.URL))

	pp, err := New(conf, testify.Logger())
	require.NoError(t, err)

	return pp, conf
}

func ask() config.Strategy {
	pass := testdata.Fake.Internet().Password()
	hash, _ := password.Hash(pass)
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

func internal() config.Strategy {
	return config.Strategy{
		Engine: config.EngineInternal,
		Name:   uuid.NewString(),
		Internal: config.Internal{
			Sqlx: sqlxconfig.Config{
				Uri: testdata.SqliteUri,
				Connection: sqlxconfig.Connection{
					MaxLifetime:  sqlxconfig.DefaultConnMaxLifetime,
					MaxIdletime:  sqlxconfig.DefaultConnMaxIdletime,
					MaxIdleCount: sqlxconfig.DefaultConnMaxIdleCount,
					MaxOpenCount: sqlxconfig.DefaultConnMaxOpenCount,
				},
			},
		},
	}
}

func external(uri string) config.Strategy {
	return config.Strategy{
		Engine: config.EngineExternal,
		Name:   uuid.NewString(),
		External: config.External{
			Uri: uri,
		},
	}
}

func httptestserver() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))

	return server
}
