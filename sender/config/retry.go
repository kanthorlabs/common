package config

import "github.com/kanthorlabs/common/validator"

type Retry struct {
	Count    int   `json:"count" yaml:"count" mapstructure:"count"`
	WaitTime int64 `json:"wait_time" yaml:"wait_time" mapstructure:"wait_time"`
}

func (conf *Retry) Validate() error {
	return validator.Validate(
		validator.NumberGreaterThanOrEqual("SENDER.CONFIG.RETRY.COUNT", conf.Count, 0),
		validator.NumberGreaterThanOrEqual("SENDER.CONFIG.RETRY.WAIT_TIME", conf.WaitTime, 100),
	)
}
