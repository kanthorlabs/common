package webhook

import "time"

var (
	DefaultKeyNs             = "epsec"
	HeaderId                 = "Webhook-Id"
	HeaderTimestamp          = "Webhook-Timestamp"
	HeaderSignature          = "Webhook-Signature"
	DefaultToleranceDuration = time.Minute * 5
	MaxKeys                  = 10
)

type Options struct {
	KeyNamespace string
}

type Option func(option *Options)

func KeyNamespace(ns string) Option {
	return func(option *Options) {
		option.KeyNamespace = ns
	}
}

type VerifyOptions struct {
	TimestampToleranceIgnore   bool
	TimestampToleranceDuration time.Duration
}

type VerifyOption func(option *VerifyOptions)

func TimestampToleranceDuration(duration time.Duration) VerifyOption {
	return func(option *VerifyOptions) {
		option.TimestampToleranceDuration = duration
	}
}

func TimestampToleranceIgnore() VerifyOption {
	return func(option *VerifyOptions) {
		option.TimestampToleranceIgnore = true
	}
}
