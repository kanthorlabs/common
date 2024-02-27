package idempotency

import "errors"

var (
	ErrAlreadyConnected = errors.New("IDEMPOTENCY.ALREADY_CONNECTED.ERROR")
	ErrNotConnected     = errors.New("IDEMPOTENCY.NOT_CONNECTED.ERROR")
	ErrConflict         = errors.New("IDEMPOTENCY.CONFLICT.ERROR")
)
