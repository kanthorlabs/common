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
	Error(args ...any)
	// Warn creates a log entry that includes a Key/WarnValue pair.
	Warn(args ...any)
	// Info creates a log entry that includes a Key/InfoValue pair.
	Info(args ...any)
	// Debug creates a log entry that includes a Key/DebugValue pair.
	Debug(args ...any)

	// Errorw creates a log entry that includes a Key/ErrorValue pair.
	Errorw(msg string, args ...any)
	// Warnw creates a log entry that includes a Key/WarnValue pair.
	Warnw(msg string, args ...any)
	// Infow creates a log entry that includes a Key/InfoValue pair.
	Infow(msg string, args ...any)
	// Debugw creates a log entry that includes a Key/DebugValue pair.
	Debugw(msg string, args ...any)

	Errorf(format string, args ...any)
	Warnf(format string, args ...any)
	Infof(format string, args ...any)
	Debugf(format string, args ...any)

	// With returns a new Logger with given args as default Key/Value pairs.
	With(args ...any) Logger
}
