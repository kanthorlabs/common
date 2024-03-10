package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/kanthorlabs/common/gatekeeper"
	gkentities "github.com/kanthorlabs/common/gatekeeper/entities"
	"github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/common/passport"
	ppentities "github.com/kanthorlabs/common/passport/entities"
)

var (
	HeaderAuthzTenant string = "X-Authorization-Tenant"
)

func Authz(authz gatekeeper.Gatekeeper, scope string) Middleware {
	if scope == "" {
		panic(errors.New("GATEWAY.AUTHZ.SCOPE_EMPTY.ERROR"))
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			acc, exist := ctx.Value(passport.CtxAccount).(*ppentities.Account)
			if !exist {
				writer.ErrUnauthorized(w, writer.ErrorString("GATEWAY.AUTHZ.ACCOUNT_EMPTY.ERROR"))
				return
			}

			tenant := parseTenant(acc, r)
			if tenant == "" {
				writer.ErrUnauthorized(w, writer.ErrorString("GATEWAY.AUTHZ.TENANT_EMPTY.ERROR"))
				return
			}

			evaluation := &gkentities.Evaluation{
				Tenant:   tenant,
				Username: acc.Username,
			}

			chictx := chi.RouteContext(r.Context())
			urlpattern := object(chictx.RoutePattern(), chictx.RoutePath)
			if urlpattern == "" {
				writer.ErrUnauthorized(w, writer.ErrorString("GATEWAY.AUTHZ.OBJECT_EMPTY.ERROR"))
				return
			}

			permission := &gkentities.Permission{
				Scope:  scope,
				Action: r.Method,
				Object: urlpattern,
			}
			err := authz.Enforce(ctx, evaluation, permission)
			if err != nil {
				writer.ErrUnauthorized(w, writer.Error(err))
				return
			}

			ctx = context.WithValue(ctx, gatekeeper.CtxTenantId, tenant)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func parseTenant(acc *ppentities.Account, r *http.Request) string {
	// prioritize the embedded tenant id inside account metadata
	if acc.Metadata != nil {
		id, has := acc.Metadata.Get(string(gatekeeper.CtxTenantId))
		if has {
			return id.(string)
		}
	}

	return r.Header.Get(HeaderAuthzTenant)
}

func object(pattern, end string) string {
	// the pattern usually ends with "/*" which is not a valid pattern to evaluate
	// so we need to replace it with a valid pattern
	return strings.Replace(pattern, "/*", "", -1) + end
}
