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
	Timeout int `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	MaxTry  int `json:"max_try" yaml:"max_try" mapstructure:"max_try"`
}

func (conf *Config) Validate() error {
	return validator.Validate(
		validator.StringRequired("HEALTHCHECK.CONFIG.DEST", conf.Dest),
		validator.NumberGreaterThanOrEqual("HEALTHCHECK.CONFIG.READINESS.TIMEOUT", conf.Readiness.Timeout, 0),
		validator.NumberGreaterThanOrEqual("HEALTHCHECK.CONFIG.READINESS.MAX_TRY", conf.Readiness.MaxTry, 0),
		validator.NumberGreaterThanOrEqual("HEALTHCHECK.CONFIG.LIVENESS.TIMEOUT", conf.Liveness.Timeout, 0),
		validator.NumberGreaterThanOrEqual("HEALTHCHECK.CONFIG.LIVENESS.MAX_TRY", conf.Liveness.MaxTry, 0),
	)
}

func Default(name string, t int) *Config {
	return &Config{
		Dest:      path.Join(os.TempDir(), name),
		Readiness: Check{Timeout: t, MaxTry: 3},
		Liveness:  Check{Timeout: t, MaxTry: 3},
	}
}
