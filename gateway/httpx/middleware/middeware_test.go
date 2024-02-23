package middleware

import "github.com/kanthorlabs/common/gateway/config"

var testconif = &config.Config{
	Addr:    ":8080",
	Timeout: 60000,
	Cors: config.Cors{
		MaxAge: 86400000,
	},
}
