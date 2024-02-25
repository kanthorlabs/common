package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/kanthorlabs/common/gateway/httpx/writer"
	"github.com/kanthorlabs/common/passport"
	"github.com/kanthorlabs/common/passport/config"
	"github.com/kanthorlabs/common/passport/entities"
)

var (
	HeaderAuthnCredentials string = "Authorization"
	HeaderAuthnEngine      string = "X-Authorization-Engine"
)

func Authn(authn passport.Passport, fallback string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			engine := r.Header.Get(HeaderAuthnEngine)
			if engine == "" {
				engine = fallback
			}

			credentials, err := parseAuthnCredentials(engine, r)
			if err != nil {
				writer.ErrUnauthorized(w, writer.Error(err))
				return
			}

			acc, err := authn.Verify(ctx, engine, credentials)
			if err != nil {
				writer.ErrUnauthorized(w, writer.Error(err))
				return
			}

			ctx = context.WithValue(ctx, passport.CtxAccount, acc)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func parseAuthnCredentials(engine string, r *http.Request) (*entities.Credentials, error) {
	if engine == config.EngineAsk || engine == config.EngineDurability {
		username, password, ok := r.BasicAuth()
		if !ok {
			return nil, errors.New("GATEWAY.AUTHN.CRENDEITALS_PARSE.ERROR")
		}
		return &entities.Credentials{Username: username, Password: password}, nil
	}

	return nil, errors.New("GATEWAY.AUTHN.ENGINE_UNKNOWN.ERROR")
}
