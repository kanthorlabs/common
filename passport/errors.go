package passport

import "errors"

var (
	ErrNotReady           = errors.New("PASSPORT.NOT_READY.ERROR")
	ErrNotLive            = errors.New("PASSPORT.NOT_LIVE.ERROR")
	ErrAlreadyConnected   = errors.New("PASSPORT.ALREADY_CONNECTED.ERROR")
	ErrNotConnected       = errors.New("PASSPORT.NOT_CONNECTED.ERROR")
	ErrStrategyNotFound   = errors.New("PASSPORT.STRATEGY.NOT_FOUND.ERROR")
	ErrStrategyDuplicated = errors.New("PASSPORT.STRATEGY.REGISTER.DUPLICATED.ERROR")
)
