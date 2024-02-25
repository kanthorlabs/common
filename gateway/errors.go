package gateway

import "errors"

var (
	ErrAlreadyStarted    = errors.New("GATEWAY.ALREADY_STARTED.ERROR")
	ErrNotStarted        = errors.New("GATEWAY.NOT_STARTED.ERROR")
	ErrHandlerNotSet     = errors.New("GATEWAY.HANDLER.NOT_SET.ERROR")
	ErrHandlerAlreadySet = errors.New("GATEWAY.HANDLER.ALREADY_SET.ERROR")
)
