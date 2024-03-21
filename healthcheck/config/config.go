package config

import (
	"os"
	"path"

	"github.com/kanthorlabs/common/validator"
)

type Config struct {
	Dest      string `json:"dest" yaml:"dest" mapstructure:"dest"`
	Readiness Check  `json:"readiness" yaml:"readiness" mapstructure:"readiness"`
	Liveness  Check  `json:"liveness" yaml:"liveness" mapstructure:"liveness"`
}

type Check struct {
	Interval int64 `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
}

func (conf *Config) Validate() error {
	return validator.Validate(
		validator.StringRequired("HEALTHCHECK.CONFIG.DEST", conf.Dest),
		validator.NumberGreaterThanOrEqual("HEALTHCHECK.CONFIG.READINESS.INTERVAL", conf.Readiness.Interval, 1000),
		validator.NumberGreaterThanOrEqual("HEALTHCHECK.CONFIG.LIVENESS.INTERVAL", conf.Liveness.Interval, 1000),
	)
}

func Default(name string, t int64) *Config {
	return &Config{
		Dest:      path.Join(os.TempDir(), name),
		Readiness: Check{Interval: t},
		Liveness:  Check{Interval: t},
	}
}
