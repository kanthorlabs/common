package testify

import "github.com/kanthorlabs/common/logging"

func Logger() logging.Logger {
	logger, _ := logging.NewNoop()
	return logger
}
