package strategies

import "errors"

var (
	ErrLogin           = errors.New("PASSPORT.ASK.LOGIN.ERROR")
	ErrRegister        = errors.New("PASSPORT.ASK.REGISTER.ERROR")
	ErrAccountNotFound = errors.New("PASSPORT.ASK.ACCOUNT.NOT_FOUND.ERROR")
	ErrDeactivate      = errors.New("PASSPORT.ASK.DEACTIVATE.ERROR")
)
