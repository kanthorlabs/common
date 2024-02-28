package logging

import "go.uber.org/zap"

type noop struct {
	*zap.SugaredLogger
}

// With returns a new no-op logger.
func (logger *noop) With(args ...any) Logger {
	return logger
}

func NewNoop() (Logger, error) {
	return &noop{
		SugaredLogger: zap.NewNop().Sugar(),
	}, nil
}
