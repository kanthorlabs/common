package config

import "github.com/kanthorlabs/common/validator"

type Config struct {
	Size  int   `json:"size" yaml:"size" mapstructure:"size"`
	Close Close `json:"close" yaml:"close" mapstructure:"close"`
	Half  Half  `json:"half" yaml:"half" mapstructure:"half"`
	Open  Open  `json:"open" yaml:"open" mapstructure:"open"`
}

func (conf *Config) Validate() error {
	err := validator.Validate(validator.NumberGreaterThan("CIRCUIT_BREAKER.CONFIG.SIZE", conf.Size, 1))
	if err != nil {
		return err
	}

	if err := conf.Close.Validate(); err != nil {
		return err
	}
	if err := conf.Half.Validate(); err != nil {
		return err
	}
	if err := conf.Open.Validate(); err != nil {
		return err
	}
	return nil
}

type Close struct {
	CleanupInterval int64 `json:"cleanup_interval" yaml:"cleanup_interval" mapstructure:"cleanup_interval"`
}

func (conf *Close) Validate() error {
	return validator.Validate(
		validator.NumberGreaterThan("CIRCUIT_BREAKER.CONFIG.CLOSE.CLEANUP_INTERVAL", conf.CleanupInterval, 1000),
	)
}

type Half struct {
	PassthroughRequests uint32 `json:"passthrough_requests" yaml:"passthrough_requests" mapstructure:"passthrough_requests"`
}

func (conf *Half) Validate() error {
	return validator.Validate(
		validator.NumberGreaterThan("CIRCUIT_BREAKER.CONFIG.HALF.PASSTHROUGH_REQUESTS", conf.PassthroughRequests, 0),
	)
}

type Open struct {
	Duration   int64          `json:"duration" yaml:"duration" mapstructure:"duration"`
	Conditions OpenConditions `json:"conditions" yaml:"conditions" mapstructure:"conditions"`
}

func (conf *Open) Validate() error {
	err := validator.Validate(
		validator.NumberGreaterThan("CIRCUIT_BREAKER.CONFIG.OPEN.DURATION", conf.Duration, 1000),
	)
	if err != nil {
		return err
	}

	if err := conf.Conditions.Validate(); err != nil {
		return err
	}

	return nil
}

type OpenConditions struct {
	ErrorConsecutive uint32  `json:"error_consecutive" yaml:"error_consecutive" mapstructure:"error_consecutive"`
	ErrorRatio       float32 `json:"error_ratio" yaml:"error_ratio" mapstructure:"error_ratio"`
}

func (conf *OpenConditions) Validate() error {
	return validator.Validate(
		validator.NumberGreaterThan("CIRCUIT_BREAKER.CONFIG.OPEN.CONDITION.ERROR_CONSECUTIVE", conf.ErrorConsecutive, 1),
		validator.NumberGreaterThan("CIRCUIT_BREAKER.CONFIG.OPEN.CONDITION.ERROR_RATIO", conf.ErrorRatio, 0.0),
		validator.NumberLessThan("CIRCUIT_BREAKER.CONFIG.OPEN.CONDITION.ERROR_RATIO", conf.ErrorRatio, 1.0),
	)
}
