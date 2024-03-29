package strategies

import "errors"

var (
	ErrNotReady           = errors.New("PASSPORT.STRATEGY.NOT_READY.ERROR")
	ErrNotLive            = errors.New("PASSPORT.STRATEGY.NOT_LIVE.ERROR")
	ErrAlreadyConnected   = errors.New("PASSPORT.STRATEGY.ALREADY_CONNECTED.ERROR")
	ErrNotConnected       = errors.New("PASSPORT.STRATEGY.NOT_CONNECTED.ERROR")
	ErrParseCredentials   = errors.New("PASSPORT.STRATEGY.PARSE_CREDENTIALS.ERROR")
	ErrCredentialsScheme  = errors.New("PASSPORT.STRATEGY.CREDENTIALS_SCHEME.ERROR")
	ErrLogin              = errors.New("PASSPORT.STRATEGY.LOGIN.ERROR")
	ErrRegister           = errors.New("PASSPORT.STRATEGY.REGISTER.ERROR")
	ErrAccountNotFound    = errors.New("PASSPORT.STRATEGY.ACCOUNT.NOT_FOUND.ERROR")
	ErrAccountDeactivated = errors.New("PASSPORT.STRATEGY.ACCOUNT.DEACTIVATED.ERROR")
	ErrDeactivate         = errors.New("PASSPORT.STRATEGY.DEACTIVATE.ERROR")
	ErrList               = errors.New("PASSPORT.STRATEGY.LIST.ERROR")
	ErrUpdate             = errors.New("PASSPORT.STRATEGY.UPDATE.ERROR")
)
