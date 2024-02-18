package gateway

import "errors"

var (
	ErrAlreadyStarted = errors.New("GATEWAY.ALREADY_STARTED.ERROR")
	ErrNotStarted     = errors.New("GATEWAY.NOT_STARTED.ERROR")
)
