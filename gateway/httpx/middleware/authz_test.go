package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	gk "github.com/kanthorlabs/common/gatekeeper"
	gkentities "github.com/kanthorlabs/common/gatekeeper/entities"
	"github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/common/mocks/gatekeeper"
	"github.com/kanthorlabs/common/passport"
	"github.com/kanthorlabs/common/passport/entities"
	"github.com/kanthorlabs/common/safe"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAuthz(t *testing.T) {
	s := chi.NewRouter()

	authz := gatekeeper.NewGatekeeper(t)
	tenantId := uuid.NewString()
	evaluation := &gkentities.Evaluation{
		Tenant:   tenantId,
		Username: account.Username,
	}
	path := "/api/workspace/{id}/extra"
	permission := &gkentities.Permission{
		Scope:  testdata.Fake.RandomStringWithLength(10),
		Action: http.MethodGet,
		Object: path,
	}

	s.Route("/api", func(sr chi.Router) {
		sr.Route("/workspace", func(ssr chi.Router) {
			ssr.Route("/{id}", func(sssr chi.Router) {
				sssr.Use(Authz(authz, permission.Scope))

				sssr.Get("/extra", func(w http.ResponseWriter, r *http.Request) {
					writer.Ok(w, writer.M{})
				})
			})
		})

	})

	t.Run("OK - tenant from header", func(st *testing.T) {
		req, err := http.NewRequest(
			http.MethodGet,
			strings.Replace(path, "{id}", uuid.NewString(), -1),
			nil,
		)
		require.NoError(st, err)

		req = req.WithContext(context.WithValue(req.Context(), passport.CtxAccount, account))
		req.Header.Set(HeaderAuthzTenant, tenantId)

		authz.EXPECT().
			Enforce(mock.Anything, evaluation, permission).
			Return(nil).
			Once()

		res := httptest.NewRecorder()
		s.ServeHTTP(res, req)

		require.Equal(st, http.StatusOK, res.Code)
	})

	t.Run("OK - tenant from metadata", func(st *testing.T) {
		req, err := http.NewRequest(
			http.MethodGet,
			strings.Replace(path, "{id}", uuid.NewString(), -1),
			nil,
		)
		require.NoError(st, err)

		// prepare attached tenant id in account metadata
		attachedTenantId := uuid.NewString()
		accountWithTenantId := &entities.Account{
			Username:     account.Username,
			PasswordHash: account.PasswordHash,
			Metadata:     &safe.Metadata{},
		}
		accountWithTenantId.Metadata.Set(string(gk.CtxTenantId), attachedTenantId)
		req = req.WithContext(context.WithValue(req.Context(), passport.CtxAccount, accountWithTenantId))
		req.Header.Set(HeaderAuthzTenant, tenantId)

		// expect we use the attached tenant id to be used to enforce permission
		eval := &gkentities.Evaluation{
			Tenant:   attachedTenantId,
			Username: account.Username,
		}
		authz.EXPECT().
			Enforce(mock.Anything, eval, permission).
			Return(nil).
			Once()

		res := httptest.NewRecorder()
		s.ServeHTTP(res, req)

		require.Equal(st, http.StatusOK, res.Code)
	})

	t.Run("KO - not permission error", func(st *testing.T) {
		req, err := http.NewRequest(http.MethodGet, path, nil)
		require.NoError(st, err)

		req = req.WithContext(context.WithValue(req.Context(), passport.CtxAccount, account))
		req.Header.Set(HeaderAuthzTenant, tenantId)

		exception := errors.New("GATEKEEPER.PERMISSION.DENINED.ERROR")
		authz.EXPECT().
			Enforce(mock.Anything, evaluation, permission).
			Return(exception).
			Once()

		res := httptest.NewRecorder()
		s.ServeHTTP(res, req)

		require.Equal(st, http.StatusUnauthorized, res.Code)

		var body writer.M
		err = json.Unmarshal(res.Body.Bytes(), &body)
		require.NoError(st, err)

		require.Contains(st, body["error"], exception.Error())
	})

	t.Run("KO - no account error", func(st *testing.T) {
		req, err := http.NewRequest(http.MethodGet, path, nil)
		require.NoError(st, err)

		req.Header.Set(HeaderAuthzTenant, tenantId)

		res := httptest.NewRecorder()
		s.ServeHTTP(res, req)

		require.Equal(st, http.StatusUnauthorized, res.Code)

		var body writer.M
		err = json.Unmarshal(res.Body.Bytes(), &body)
		require.NoError(st, err)

		require.Contains(st, body["error"], "GATEWAY.AUTHZ.ACCOUNT_EMPTY.ERROR")
	})

	t.Run("KO - no tenant error", func(st *testing.T) {
		req, err := http.NewRequest(http.MethodGet, path, nil)
		require.NoError(st, err)

		req = req.WithContext(context.WithValue(req.Context(), passport.CtxAccount, account))

		res := httptest.NewRecorder()
		s.ServeHTTP(res, req)

		require.Equal(st, http.StatusUnauthorized, res.Code)

		var body writer.M
		err = json.Unmarshal(res.Body.Bytes(), &body)
		require.NoError(st, err)

		require.Contains(st, body["error"], "GATEWAY.AUTHZ.TENANT_EMPTY.ERROR")
	})

	t.Run("KO - no route pattern", func(st *testing.T) {
		s := chi.NewRouter()
		authz := gatekeeper.NewGatekeeper(t)
		// top level will not work because we cannot detect the mattching pattern
		s.Use(Authz(authz, testdata.Fake.RandomStringWithLength(10)))
		s.Get("/undetectable", func(w http.ResponseWriter, r *http.Request) {
			writer.Ok(w, writer.M{})
		})

		req, err := http.NewRequest(http.MethodGet, "/undetectable", nil)
		require.NoError(st, err)

		req = req.WithContext(context.WithValue(req.Context(), passport.CtxAccount, account))
		req.Header.Set(HeaderAuthzTenant, tenantId)

		res := httptest.NewRecorder()
		s.ServeHTTP(res, req)

		require.Equal(st, http.StatusUnauthorized, res.Code)

		var body writer.M
		err = json.Unmarshal(res.Body.Bytes(), &body)
		require.NoError(st, err)

		require.Contains(st, body["error"], "GATEWAY.AUTHZ.OBJECT_EMPTY.ERROR")
	})
}

func TestAuthz_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			require.ErrorContains(t, r.(error), "GATEWAY.AUTHZ.SCOPE_EMPTY.ERROR")
		}
	}()

	s := chi.NewRouter()
	authz := gatekeeper.NewGatekeeper(t)
	s.Use(Authz(authz, ""))
}
