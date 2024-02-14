package passport

import "errors"

var (
	ErrStrategyNotFound   = errors.New("PASSPORT.STRATEGY.NOT_FOUND.ERROR")
	ErrStrategyDuplicated = errors.New("PASSPORT.STRATEGY.REGISTER.DUPLICATED.ERROR")
)
