package logging

import (
	"github.com/kanthorlabs/common/configuration"
	"github.com/kanthorlabs/common/logging/config"
)

func New(provider configuration.Provider) (Logger, error) {
	conf, err := config.New(provider)
	if err != nil {
		return nil, err
	}
	return NewZap(conf)
}

// Logger is the logger abstraction. It largely follows zap structure.
type Logger interface {
	// Error creates a log entry that includes a Key/ErrorValue pair.
	Error(args ...interface{})
	// Warn creates a log entry that includes a Key/WarnValue pair.
	Warn(args ...interface{})
	// Info creates a log entry that includes a Key/InfoValue pair.
	Info(args ...interface{})
	// Debug creates a log entry that includes a Key/DebugValue pair.
	Debug(args ...interface{})

	// Errorw creates a log entry that includes a Key/ErrorValue pair.
	Errorw(msg string, args ...interface{})
	// Warnw creates a log entry that includes a Key/WarnValue pair.
	Warnw(msg string, args ...interface{})
	// Infow creates a log entry that includes a Key/InfoValue pair.
	Infow(msg string, args ...interface{})
	// Debugw creates a log entry that includes a Key/DebugValue pair.
	Debugw(msg string, args ...interface{})

	Errorf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})

	// With returns a new Logger with given args as default Key/Value pairs.
	With(args ...interface{}) Logger
}
