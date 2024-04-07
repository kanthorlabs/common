package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/kanthorlabs/common/cache"
	"github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/common/passport"
	ppentities "github.com/kanthorlabs/common/passport/entities"
)

var (
	HeaderAuthnCredentials string = "Authorization"
	HeaderAuthnStrategy    string = "X-Authorization-Stategy"
)

func Authn(authn passport.Passport, options ...AuthnOption) Middleware {
	conf := &AuthnConfig{
		Cache:     cache.NewNoop(),
		Fallback:  "",
		ExpiresIn: time.Hour,
	}
	for _, with := range options {
		with(conf)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			name := r.Header.Get(HeaderAuthnStrategy)
			if name == "" {
				name = conf.Fallback
			}
			tokens := ppentities.Tokens{Access: r.Header.Get(HeaderAuthnCredentials)}

			acc, err := cache.GetOrSet(
				conf.Cache,
				ctx,
				cache.EncodeKey(name, tokens.Access),
				conf.ExpiresIn,
				func() (*ppentities.Account, error) {
					strategy, err := authn.Strategy(name)
					if err != nil {
						return nil, err
					}

					return strategy.Verify(ctx, tokens)
				},
			)
			if err != nil {
				writer.ErrUnauthorized(w, writer.Error(err))
				return
			}

			ctx = context.WithValue(ctx, passport.CtxAccount, acc)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type AuthnConfig struct {
	Cache     cache.Cache
	Fallback  string
	ExpiresIn time.Duration
}

type AuthnOption func(*AuthnConfig) error

func AuthnWithCache(c cache.Cache) AuthnOption {
	return func(o *AuthnConfig) error {
		o.Cache = c
		return nil
	}
}

func AuthnWithFallback(fallback string) AuthnOption {
	return func(o *AuthnConfig) error {
		o.Fallback = fallback
		return nil
	}
}

func AuthnWithExpiresIn(d time.Duration) AuthnOption {
	return func(o *AuthnConfig) error {
		o.ExpiresIn = d
		return nil
	}
}
