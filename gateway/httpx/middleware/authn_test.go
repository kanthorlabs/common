package middleware

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kanthorlabs/common/cipher/password"
	"github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/common/mocks/passport"
	"github.com/kanthorlabs/common/passport/config"
	"github.com/kanthorlabs/common/passport/entities"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	user        = uuid.NewString()
	pass        = uuid.NewString()
	hash, _     = password.HashString(pass)
	credentials = base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))
)

func TestAuthn(t *testing.T) {
	s := chi.NewRouter()
	authn := passport.NewPassport(t)
	s.Use(Authn(authn, config.EngineAsk))

	path := "/"
	s.Get(path, func(w http.ResponseWriter, r *http.Request) {
		writer.Ok(w, writer.M{})
	})

	t.Run("OK - fallback", func(st *testing.T) {
		req, err := http.NewRequest(http.MethodGet, path, nil)
		req.Header.Set(HeaderAuthnCredentials, "basic "+credentials)
		require.Nil(st, err)

		// only care about the engine
		authn.EXPECT().
			Verify(mock.Anything, config.EngineAsk, mock.Anything).
			Return(&entities.Account{Username: user, PasswordHash: hash}, nil).
			Once()

		res := httptest.NewRecorder()
		s.ServeHTTP(res, req)

		require.Equal(st, http.StatusOK, res.Code)
	})

	t.Run("OK - set via header", func(st *testing.T) {
		req, err := http.NewRequest(http.MethodGet, path, nil)
		req.Header.Set(HeaderAuthnEngine, config.EngineDurability)
		req.Header.Set(HeaderAuthnCredentials, "basic "+credentials)
		require.Nil(st, err)

		// only care about the engine
		authn.EXPECT().
			Verify(mock.Anything, config.EngineDurability, &entities.Credentials{Username: user, Password: pass}).
			Return(&entities.Account{Username: user, PasswordHash: hash}, nil).
			Once()

		res := httptest.NewRecorder()
		s.ServeHTTP(res, req)

		require.Equal(st, http.StatusOK, res.Code)
	})

	t.Run("KO - unknown engine", func(st *testing.T) {
		req, err := http.NewRequest(http.MethodGet, path, nil)
		req.Header.Set(HeaderAuthnEngine, testdata.Fake.Blood().Name())
		req.Header.Set(HeaderAuthnCredentials, "basic "+credentials)
		require.Nil(st, err)

		res := httptest.NewRecorder()
		s.ServeHTTP(res, req)

		require.Equal(st, http.StatusUnauthorized, res.Code)

		var body writer.M
		err = json.Unmarshal(res.Body.Bytes(), &body)
		require.Nil(st, err)

		require.Contains(st, body["error"], "GATEWAY.AUTHN.ENGINE_UNKNOWN.ERROR")
	})

	t.Run("KO - parse credentials error", func(st *testing.T) {
		req, err := http.NewRequest(http.MethodGet, path, nil)
		require.Nil(st, err)

		res := httptest.NewRecorder()
		s.ServeHTTP(res, req)

		require.Equal(st, http.StatusUnauthorized, res.Code)

		var body writer.M
		err = json.Unmarshal(res.Body.Bytes(), &body)
		require.Nil(st, err)

		require.Contains(st, body["error"], "GATEWAY.AUTHN.CRENDEITALS_PARSE.ERROR")
	})

	t.Run("KO - verify error", func(st *testing.T) {
		req, err := http.NewRequest(http.MethodGet, path, nil)
		req.Header.Set(HeaderAuthnCredentials, "basic "+credentials)
		require.Nil(st, err)

		expected := errors.New("PASSPORT.ACCOUNT_NOT_FOUND.ERROR")
		authn.EXPECT().
			Verify(mock.Anything, config.EngineAsk, &entities.Credentials{Username: user, Password: pass}).
			Return(nil, expected).
			Once()

		res := httptest.NewRecorder()
		s.ServeHTTP(res, req)

		require.Equal(st, http.StatusUnauthorized, res.Code)

		var body writer.M
		err = json.Unmarshal(res.Body.Bytes(), &body)
		require.Nil(st, err)

		require.Contains(st, body["error"], expected.Error())
	})

}
