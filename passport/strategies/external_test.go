package strategies

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/passport/config"
	"github.com/kanthorlabs/common/passport/entities"
	"github.com/kanthorlabs/common/project"
	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/testify"
	"github.com/kanthorlabs/common/utils"
	"github.com/stretchr/testify/require"
)

var (
	testoktext       = "###OK###"
	testkostatus     = http.StatusInternalServerError
	testkostatustext = http.StatusText(testkostatus)
)

func TestExternal_New(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		conf := &config.External{
			Uri: fmt.Sprintf("http://localhost:%d", testdata.Fake.IntBetween(2000, 10000)),
		}
		_, err := NewExternal(conf, testify.Logger())
		require.NoError(st, err)
	})

	t.Run("KO - configuration error", func(st *testing.T) {
		_, err := NewExternal(&config.External{}, testify.Logger())
		require.ErrorContains(st, err, "PASSPORT.CONFIG.EXTERNAL")
	})
}

func TestExternal_Connect(t *testing.T) {
	conf, server := externalsetup(nil)
	defer server.Close()

	t.Run("OK", func(st *testing.T) {
		strategy, err := NewExternal(conf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, strategy.Connect(context.Background()))
	})

	t.Run("KO - already connected error", func(st *testing.T) {
		strategy, err := NewExternal(conf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, strategy.Connect(context.Background()))
		require.ErrorIs(st, strategy.Connect(context.Background()), ErrAlreadyConnected)
	})
}

func TestExternal_Readiness(t *testing.T) {
	conf, server := externalsetup(nil)
	defer server.Close()

	t.Run("OK", func(st *testing.T) {
		strategy, err := NewExternal(conf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, strategy.Connect(context.Background()))
		require.NoError(st, strategy.Readiness())
	})

	t.Run("OK - disconnected", func(st *testing.T) {
		strategy, err := NewExternal(conf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, strategy.Connect(context.Background()))
		require.NoError(st, strategy.Disconnect(context.Background()))
		require.NoError(st, strategy.Readiness())
	})

	t.Run("KO - not connected error", func(st *testing.T) {
		strategy, err := NewExternal(conf, testify.Logger())
		require.NoError(st, err)

		require.ErrorIs(st, strategy.Readiness(), ErrNotConnected)
	})
}

func TestExternal_Liveness(t *testing.T) {
	conf, server := externalsetup(nil)
	defer server.Close()

	t.Run("OK", func(st *testing.T) {
		strategy, err := NewExternal(conf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, strategy.Connect(context.Background()))
		require.NoError(st, strategy.Liveness())
	})

	t.Run("OK - disconnected", func(st *testing.T) {
		strategy, err := NewExternal(conf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, strategy.Connect(context.Background()))
		require.NoError(st, strategy.Disconnect(context.Background()))
		require.NoError(st, strategy.Liveness())
	})

	t.Run("KO - not connected error", func(st *testing.T) {
		strategy, err := NewExternal(conf, testify.Logger())
		require.NoError(st, err)

		require.ErrorIs(st, strategy.Liveness(), ErrNotConnected)
	})
}

func TestExternal_Disconnect(t *testing.T) {
	conf, server := externalsetup(nil)
	defer server.Close()

	t.Run("OK", func(st *testing.T) {
		strategy, err := NewExternal(conf, testify.Logger())
		require.NoError(st, err)

		require.NoError(st, strategy.Connect(context.Background()))
		require.NoError(st, strategy.Disconnect(context.Background()))
	})

	t.Run("KO - not connected error", func(st *testing.T) {
		strategy, err := NewExternal(conf, testify.Logger())
		require.NoError(st, err)

		require.ErrorIs(st, strategy.Disconnect(context.Background()), ErrNotConnected)
	})
}

func TestExternal_Register(t *testing.T) {
	conf, server := externalsetup(nil)
	defer server.Close()

	strategy, err := NewExternal(conf, testify.Logger())
	require.NoError(t, err)

	strategy.Connect(context.Background())
	defer strategy.Disconnect(context.Background())

	t.Run("KO", func(st *testing.T) {
		acc := entities.Account{
			Username:  uuid.NewString(),
			Password:  uuid.NewString(),
			Name:      testdata.Fake.Internet().User(),
			CreatedAt: time.Now().UnixMilli(),
			UpdatedAt: time.Now().UnixMilli(),
		}

		require.ErrorContains(st, strategy.Register(context.Background(), acc), "UNIMPLEMENT.ERROR")
	})
}

func TestExternal_Login(t *testing.T) {
	conf, server := externalsetup(nil)
	defer server.Close()

	strategy, err := NewExternal(conf, testify.Logger())
	require.NoError(t, err)

	strategy.Connect(context.Background())
	defer strategy.Disconnect(context.Background())

	t.Run("KO", func(st *testing.T) {
		creds := entities.Credentials{
			Region:   project.Region(),
			Username: uuid.NewString(),
			Password: testoktext + uuid.NewString(),
		}

		_, err := strategy.Login(context.Background(), creds)
		require.ErrorContains(st, err, "UNIMPLEMENT.ERROR")
	})
}

func TestExternal_Logout(t *testing.T) {
	conf, server := externalsetup(nil)
	defer server.Close()

	strategy, err := NewExternal(conf, testify.Logger())
	require.NoError(t, err)

	strategy.Connect(context.Background())
	defer strategy.Disconnect(context.Background())

	t.Run("KO", func(st *testing.T) {
		tokens := entities.Tokens{
			Access: testdata.Fake.RandomStringWithLength(256),
		}

		require.ErrorContains(st, strategy.Logout(context.Background(), tokens), "UNIMPLEMENT.ERROR")
	})
}

func TestExternal_Verify(t *testing.T) {
	var emptyres = "###EMPTY###"

	conf, server := externalsetup(func(w http.ResponseWriter, r *http.Request) {
		account := &entities.Account{
			Username:  uuid.NewString(),
			Password:  uuid.NewString(),
			Name:      testdata.Fake.Internet().User(),
			CreatedAt: time.Now().UnixMilli(),
			UpdatedAt: time.Now().UnixMilli(),
		}

		if strings.Contains(r.Context().Value(testctxkeydata).(string), emptyres) {
			account.Username = ""
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(utils.Stringify(account)))
	})
	defer server.Close()

	strategy, err := NewExternal(conf, testify.Logger())
	require.NoError(t, err)

	strategy.Connect(context.Background())
	defer strategy.Disconnect(context.Background())

	t.Run("OK", func(st *testing.T) {
		tokens := entities.Tokens{
			Access: testoktext + testdata.Fake.RandomStringWithLength(256),
		}

		acc, err := strategy.Verify(context.Background(), tokens)
		require.NoError(st, err)
		require.NotNil(st, acc)
	})

	t.Run("KO - external error", func(st *testing.T) {
		tokens := entities.Tokens{
			Access: testdata.Fake.RandomStringWithLength(256),
		}

		acc, err := strategy.Verify(context.Background(), tokens)
		require.ErrorContains(st, err, testkostatustext)
		require.Nil(st, acc)
	})

	t.Run("KO - returning account error", func(st *testing.T) {
		tokens := entities.Tokens{
			Access: testoktext + testdata.Fake.RandomStringWithLength(256) + emptyres,
		}

		acc, err := strategy.Verify(context.Background(), tokens)
		require.ErrorContains(st, err, "PASSPORT.ACCOUNT.")
		require.Nil(st, acc)
	})
}

func TestExternalManagement(t *testing.T) {
	conf, server := externalsetup(nil)
	defer server.Close()

	strategy, err := NewExternal(conf, testify.Logger())
	require.NoError(t, err)

	t.Run("KO", func(st *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				require.ErrorIs(t, r.(error), ErrNotConnected)
			}
		}()
		strategy.Management()
	})
}

func TestExternalManagement_Deactivate(t *testing.T) {
	conf, server := externalsetup(nil)
	defer server.Close()

	strategy, err := NewExternal(conf, testify.Logger())
	require.NoError(t, err)

	strategy.Connect(context.Background())
	defer strategy.Disconnect(context.Background())

	t.Run("KO - unimplement error", func(st *testing.T) {
		err := strategy.Management().Deactivate(context.Background(), uuid.NewString(), time.Now().UnixMilli())
		require.ErrorContains(st, err, "UNIMPLEMENT.ERROR")
	})
}

func TestExternalManagement_List(t *testing.T) {
	conf, server := externalsetup(nil)
	defer server.Close()

	strategy, err := NewExternal(conf, testify.Logger())
	require.NoError(t, err)

	strategy.Connect(context.Background())
	defer strategy.Disconnect(context.Background())

	t.Run("KO - unimplement error", func(st *testing.T) {
		_, err := strategy.Management().List(context.Background(), []string{})
		require.ErrorContains(st, err, "UNIMPLEMENT.ERROR")
	})
}

func TestExternalManagement_Update(t *testing.T) {
	conf, server := externalsetup(nil)
	defer server.Close()

	strategy, err := NewExternal(conf, testify.Logger())
	require.NoError(t, err)

	strategy.Connect(context.Background())
	defer strategy.Disconnect(context.Background())

	t.Run("KO - unimplement error", func(st *testing.T) {
		err := strategy.Management().Update(context.Background(), entities.Account{})
		require.ErrorContains(st, err, "UNIMPLEMENT.ERROR")
	})
}

type testctxkey string

var testctxkeydata testctxkey = "data"

func externalsetup(okhandler func(w http.ResponseWriter, r *http.Request)) (*config.External, *httptest.Server) {
	if okhandler == nil {
		okhandler = externalhandlerok
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "healthz") {
			okhandler(w, r)
			return
		}

		authn := r.Header.Get("authorization")
		if strings.Contains(authn, testoktext) {
			r = r.WithContext(context.WithValue(r.Context(), testctxkeydata, authn))

			okhandler(w, r)
			return
		}

		body, _ := io.ReadAll(r.Body)
		if strings.Contains(string(body), testoktext) {
			r = r.WithContext(context.WithValue(r.Context(), testctxkeydata, string(body)))
			okhandler(w, r)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(testkostatus)
		w.Write([]byte("oops, something went wrong"))
	}))

	conf := &config.External{
		Uri: server.URL,
	}
	return conf, server
}

var externalhandlerok = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}
