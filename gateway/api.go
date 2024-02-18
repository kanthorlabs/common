package gateway

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/kanthorlabs/common/gateway/config"
	httpx "github.com/kanthorlabs/common/gateway/httpx/middleware"
	"github.com/kanthorlabs/common/logging"
)

func NewHttpx(conf *config.Config, logger logging.Logger) (*chi.Mux, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	r := chi.NewRouter()
	r.Use(httpx.RequestId())
	r.Use(httpx.Logger(logger))
	r.Use(httpx.Recover())
	r.Use(middleware.Timeout(time.Millisecond * time.Duration(conf.Timeout)))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   conf.Cors.AllowedOrigins,
		AllowedMethods:   conf.Cors.AllowedMethods,
		AllowedHeaders:   conf.Cors.AllowedHeaders,
		ExposedHeaders:   conf.Cors.ExposedHeaders,
		AllowCredentials: conf.Cors.AllowCredentials,
		MaxAge:           conf.Cors.MaxAge,
	}))

	return r, nil
}
