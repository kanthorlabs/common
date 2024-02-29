package config

import (
	"net/url"
	"strings"

	"github.com/kanthorlabs/common/validator"
)

var EngineNats = "nats"

type Config struct {
	Name       string     `json:"name" yaml:"name" mapstructure:"name"`
	Uri        string     `json:"uri" yaml:"uri" mapstructure:"uri"`
	Nats       Nats       `json:"nats" yaml:"nats" mapstructure:"nats"`
	Publisher  Publisher  `json:"publisher" yaml:"publisher" mapstructure:"publisher"`
	Subscriber Subscriber `json:"subscriber" yaml:"subscriber" mapstructure:"subscriber"`
}

func (conf *Config) Validate() error {
	err := validator.Validate(
		validator.StringAlphaNumericUnderscore("STREAMING.CONFIG.NAME", conf.Name),
		validator.StringUri("STREAMING.CONFIG.URI", conf.Uri),
		validator.StringStartsWithOneOf("STREAMING.CONFIG.URI", conf.Uri, []string{EngineNats}),
	)
	if err != nil {
		return err
	}

	// StringUri already validated it
	uri, _ := url.ParseRequestURI(conf.Uri)
	if strings.HasPrefix(uri.Scheme, "nats") {
		if err := conf.Nats.Validate(); err != nil {
			return err
		}
	}

	if err := conf.Publisher.Validate(); err != nil {
		return err
	}

	if err := conf.Subscriber.Validate(); err != nil {
		return err
	}

	return nil
}

type Publisher struct {
	RateLimit int `json:"rate_limit" yaml:"rate_limit" mapstructure:"rate_limit"`
}

func (conf *Publisher) Validate() error {
	return validator.Validate(
		validator.NumberGreaterThan("STREAMING.CONFIG.PUBLISHER.RATE_LIMIT", conf.RateLimit, 0),
	)
}

type Subscriber struct {
	Timeout int64 `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	// MaxRetry is how many times we should try to re-deliver message if we get any error
	MaxRetry    int `json:"max_retry" yaml:"max_retry" mapstructure:"max_retry"`
	Concurrency int `json:"concurrency" yaml:"concurrency" mapstructure:"concurrency"`
}

func (conf *Subscriber) Validate() error {
	return validator.Validate(
		validator.NumberGreaterThanOrEqual("STREAMING.CONFIG.SUBSCRIBER.TIMEOUT", conf.Timeout, 1000),
		validator.NumberGreaterThanOrEqual("STREAMING.CONFIG.SUBSCRIBER.MAX_RETRY", conf.MaxRetry, 1),
		validator.NumberGreaterThan("STREAMING.CONFIG.SUBSCRIBER.CONCURRENCY", conf.Concurrency, 0),
	)
}
