package config

import (
	_ "embed"

	"github.com/kanthorlabs/common/persistence/sqlx/config"
	"github.com/kanthorlabs/common/validator"
)

var EngineRBAC = "rbac"

type Config struct {
	Engine      string      `json:"engine" yaml:"engine" mapstructure:"engine"`
	Definitions Definitions `json:"definitions" yaml:"definitions" mapstructure:"definitions"`
	Privilege   Privilege   `json:"privilege" yaml:"privilege" mapstructure:"privilege"`
}

func (conf *Config) Validate() error {
	err := validator.Validate(
		validator.StringOneOf("GATEKEEPER.CONFIG.ENGINE", conf.Engine, []string{EngineRBAC}),
	)
	if err != nil {
		return err
	}

	if err := conf.Privilege.Validate(); err != nil {
		return err
	}

	if err := conf.Definitions.Validate(); err != nil {
		return err
	}

	return nil
}

type Definitions struct {
	Uri string `json:"uri" yaml:"uri" mapstructure:"uri"`
}

func (conf *Definitions) Validate() error {
	return validator.Validate(
		validator.StringRequired("GATEKEEPER.CONFIG.DEFINITIONS.URI", conf.Uri),
		validator.StringStartsWithOneOf("GATEKEEPER.CONFIG.DEFINITIONS.URI", conf.Uri, []string{"file", "base64"}),
	)
}

type Privilege struct {
	Sqlx config.Config `json:"sqlx" yaml:"sqlx" mapstructure:"sqlx"`
}

func (conf *Privilege) Validate() error {
	return conf.Sqlx.Validate()
}
