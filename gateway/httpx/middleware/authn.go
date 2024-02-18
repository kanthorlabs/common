package middleware

import (
	"context"
	"net/http"

	"github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/common/passport"
	"github.com/kanthorlabs/common/passport/entities"
)

var (
	HeaderAuthnCredentials string = "Authorization"
	HeaderAuthnEngine      string = "X-Authorization-Engine"
	CtxAccount             ctxkey = "gateway.account"
)

func Authn(engine passport.Passport, fallback string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			username, password, ok := r.BasicAuth()
			if !ok {
				writer.ErrUnauthorized(w, writer.ErrorString("GATEWAY.AUTHN.CREDENTIAL.ERROR"))
				return
			}

			name := r.Header.Get(HeaderAuthnEngine)
			if name == "" {
				name = fallback
			}

			credentials := &entities.Credentials{
				Username: username,
				Password: password,
			}
			acc, err := engine.Verify(ctx, name, credentials)
			if err != nil {
				writer.ErrUnauthorized(w, writer.Error(err))
				return
			}

			ctx = context.WithValue(ctx, CtxAccount, acc)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
