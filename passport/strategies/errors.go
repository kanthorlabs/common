package strategies

import "errors"

var (
	ErrNotReady         = errors.New("PASSPORT.STRATEGY.NOT_READY.ERROR")
	ErrNotLive          = errors.New("PASSPORT.STRATEGY.NOT_LIVE.ERROR")
	ErrAlreadyConnected = errors.New("PASSPORT.STRATEGY.ALREADY_CONNECTED.ERROR")
	ErrNotConnected     = errors.New("PASSPORT.STRATEGY.NOT_CONNECTED.ERROR")
	ErrLogin            = errors.New("PASSPORT.STRATEGY.LOGIN.ERROR")
	ErrRegister         = errors.New("PASSPORT.STRATEGY.REGISTER.ERROR")
	ErrAccountNotFound  = errors.New("PASSPORT.STRATEGY.ACCOUNT.NOT_FOUND.ERROR")
	ErrDeactivate       = errors.New("PASSPORT.STRATEGY.DEACTIVATE.ERROR")
)
