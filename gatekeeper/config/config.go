package config

import (
	_ "embed"

	"github.com/kanthorlabs/common/persistence/sqlx/config"
	"github.com/kanthorlabs/common/validator"
)

type Config struct {
	Policy    Policy    `json:"policy" yaml:"policy" mapstructure:"policy"`
	Privilege Privilege `json:"storage" yaml:"storage" mapstructure:"storage"`
}

func (conf *Config) Validate() error {
	err := validator.Validate()
	if err != nil {
		return err
	}

	if err := conf.Policy.Validate(); err != nil {
		return err
	}

	if err := conf.Privilege.Validate(); err != nil {
		return err
	}

	return nil
}

type Policy struct {
	JudgeUri      string `json:"judge_uri" yaml:"judge_uri" mapstructure:"judge_uri"`
	PermissionUri string `json:"permission_url" yaml:"permission_url" mapstructure:"permission_url"`
}

func (conf *Policy) Validate() error {
	return validator.Validate(
		validator.StringUri("GATEKEEPER.CONFIG.POLICY.JUDGE_URI", conf.JudgeUri),
		validator.StringUri("GATEKEEPER.CONFIG.POLICY.PERMISSION_URI", conf.PermissionUri),
	)
}

type Privilege struct {
	Sqlx config.Config `json:"sqlx" yaml:"sqlx" mapstructure:"sqlx"`
}

func (conf *Privilege) Validate() error {
	return conf.Sqlx.Validate()
}
