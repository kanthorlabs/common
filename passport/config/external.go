package config

import "github.com/kanthorlabs/common/validator"

type External struct {
	Uri string `json:"uri" yaml:"uri" mapstructure:"uri"`
}

func (conf *External) Validate() error {
	return validator.Validate(
		validator.StringUri("PASSPORT.CONFIG.EXTERNAL.URI", conf.Uri),
		validator.StringStartsWithOneOf("PASSPORT.CONFIG.EXTERNAL.URI", conf.Uri, []string{"http", "https"}),
	)
}
