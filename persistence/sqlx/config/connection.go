package config

import "github.com/kanthorlabs/common/validator"

type Connection struct {
	MaxLifetime  int64 `json:"max_lifetime" yaml:"max_lifetime" mapstructure:"max_lifetime"`
	MaxIdletime  int64 `json:"max_idletime" yaml:"max_idletime" mapstructure:"max_idletime"`
	MaxIdleCount int   `json:"max_idle_count" yaml:"max_idle_count" mapstructure:"max_idle_count"`
	MaxOpenCount int   `json:"max_open_count" yaml:"max_open_count" mapstructure:"max_open_count"`
}

func (conf *Connection) Validate() error {
	return validator.Validate(
		validator.NumberInRange("SQLX.CONFIG.CONNECTION.MAX_LIFETIME", conf.MaxLifetime, DefaultConnMaxLifetime, 3600000),
		validator.NumberInRange("SQLX.CONFIG.CONNECTION.MAX_IDLETIME", conf.MaxIdletime, DefaultConnMaxIdletime, 3600000),
		validator.NumberGreaterThan("SQLX.CONFIG.CONNECTION.MAX_IDLE_COUNT", conf.MaxIdleCount, 0),
		validator.NumberGreaterThan("SQLX.CONFIG.CONNECTION.MAX_OPEN_COUNT", conf.MaxOpenCount, 0),
	)
}
