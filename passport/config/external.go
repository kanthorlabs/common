package config

import "github.com/kanthorlabs/common/validator"

type External struct {
	Uri string `json:"uri" yaml:"uri" mapstructure:"uri"`
}

func (conf *External) Validate() error {
	return validator.Validate(
		validator.StringUri("SQLX.CONFIG.URI", conf.Uri),
		validator.StringStartsWithOneOf("SQLX.CONFIG.URI", conf.Uri, []string{"http", "https"}),
	)
}
