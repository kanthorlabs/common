package sqlx

import (
	"errors"
)

var (
	ErrNotReady         = errors.New("SQLX.NOT_READY.ERROR")
	ErrNotLive          = errors.New("SQLX.NOT_LIVE.ERROR")
	ErrAlreadyConnected = errors.New("SQLX.ALREADY_CONNECTED.ERROR")
	ErrNotConnected     = errors.New("SQLX.NOT_CONNECTED.ERROR")
	ErrRecordNotFound   = errors.New("SQLX.RECORD.NOT_FOUND.ERROR")
)
