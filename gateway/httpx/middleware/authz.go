package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kanthorlabs/common/gatekeeper"
	gkEnt "github.com/kanthorlabs/common/gatekeeper/entities"
	"github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/common/passport/entities"
	ppEnt "github.com/kanthorlabs/common/passport/entities"
)

var (
	HeaderAuthzTenant string = "X-Authorization-Tenant"
	CtxTenantId       ctxkey = "gateway.tenant.id"
)

func Authz(engine gatekeeper.Gatekeeper) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			acc, exist := ctx.Value(CtxAccount).(*ppEnt.Account)
			if !exist {
				writer.ErrUnauthorized(w, writer.ErrorString("GATEWAY.AUTHZ.ACCOUNT.EMPTY.ERROR"))
				return
			}

			tenant := tenantId(acc, r.Header)
			if tenant == "" {
				writer.ErrUnauthorized(w, writer.ErrorString("GATEWAY.AUTHZ.TENANT.EMPTY.ERROR"))
				return
			}

			patterns := chi.RouteContext(ctx).RoutePatterns
			if len(patterns) == 0 {
				writer.ErrUnauthorized(w, writer.ErrorString("GATEWAY.AUTHZ.OBJECT.EMPTY.ERROR"))
			}

			for i := range patterns {
				evaluation := &gkEnt.Evaluation{
					Tenant:   tenant,
					Username: acc.Username,
				}
				permission := &gkEnt.Permission{
					Action: r.Method,
					Object: patterns[i],
				}
				err := engine.Enforce(ctx, evaluation, permission)
				if err != nil {
					writer.ErrUnauthorized(w, writer.Error(err))
					return
				}
			}

			ctx = context.WithValue(ctx, CtxTenantId, tenant)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func tenantId(acc *entities.Account, headers http.Header) string {
	// prioritize the embedded tenant id inside account metadata
	if acc.Metadata != nil {
		id, has := acc.Metadata.Get(string(CtxTenantId))
		if has {
			return id.(string)
		}
	}

	return headers.Get(HeaderAuthzTenant)
}
