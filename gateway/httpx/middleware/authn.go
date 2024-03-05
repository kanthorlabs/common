package middleware

import (
	"context"
	"net/http"

	"github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/common/passport"
)

var (
	HeaderAuthnCredentials string = "Authorization"
	HeaderAuthnStrategy    string = "X-Authorization-Stategy"
)

func Authn(authn passport.Passport, fallback string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			name := r.Header.Get(HeaderAuthnStrategy)
			if name == "" {
				name = fallback
			}

			strategy, err := authn.Strategy(name)
			if err != nil {
				writer.ErrUnauthorized(w, writer.Error(err))
				return
			}

			credentials, err := strategy.ParseCredentials(ctx, r.Header.Get(HeaderAuthnCredentials))
			if err != nil {
				writer.ErrUnauthorized(w, writer.Error(err))
				return
			}

			acc, err := strategy.Verify(ctx, credentials)
			if err != nil {
				writer.ErrUnauthorized(w, writer.Error(err))
				return
			}

			ctx = context.WithValue(ctx, passport.CtxAccount, acc)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
