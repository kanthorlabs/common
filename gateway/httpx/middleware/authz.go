package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kanthorlabs/common/gatekeeper"
	gkEnt "github.com/kanthorlabs/common/gatekeeper/entities"
	"github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/common/passport"
	"github.com/kanthorlabs/common/passport/entities"
)

var (
	HeaderAuthzTenant string = "X-Authorization-Tenant"
)

func Authz(authz gatekeeper.Gatekeeper, scope string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			acc, exist := ctx.Value(passport.CtxAccount).(*entities.Account)
			if !exist {
				writer.ErrUnauthorized(w, writer.ErrorString("GATEWAY.AUTHZ.ACCOUNT_EMPTY.ERROR"))
				return
			}

			tenant := parseTenant(acc, r)
			if tenant == "" {
				writer.ErrUnauthorized(w, writer.ErrorString("GATEWAY.AUTHZ.TENANT_EMPTY.ERROR"))
				return
			}

			patterns := chi.RouteContext(ctx).RoutePatterns
			if len(patterns) == 0 {
				writer.ErrUnauthorized(w, writer.ErrorString("GATEWAY.AUTHZ.OBJECT_EMPTY.ERROR"))
				return
			}

			for i := range patterns {
				evaluation := &gkEnt.Evaluation{
					Tenant:   tenant,
					Username: acc.Username,
				}
				permission := &gkEnt.Permission{
					Action: r.Method,
					Object: object(scope, patterns[i]),
				}
				err := authz.Enforce(ctx, evaluation, permission)
				if err != nil {
					writer.ErrUnauthorized(w, writer.Error(err))
					return
				}
			}

			ctx = context.WithValue(ctx, gatekeeper.CtxTenantId, tenant)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func parseTenant(acc *entities.Account, r *http.Request) string {
	// prioritize the embedded tenant id inside account metadata
	if acc.Metadata != nil {
		id, has := acc.Metadata.Get(string(gatekeeper.CtxTenantId))
		if has {
			return id.(string)
		}
	}

	return r.Header.Get(HeaderAuthzTenant)
}

func object(scope, pattern string) string {
	if scope == "" {
		return pattern
	}

	return fmt.Sprintf("%s::%s", scope, pattern)
}
