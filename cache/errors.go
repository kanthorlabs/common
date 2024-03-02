package cache

import "errors"

var (
	ErrAlreadyConnected = errors.New("CACHE.ALREADY_CONNECTED.ERROR")
	ErrNotConnected     = errors.New("CACHE.NOT_CONNECTED.ERROR")
	ErrEntryNotFound    = errors.New("CACHE.ENTRY.NOT_FOUND.ERROR")
	ErrKeyEmpty         = errors.New("CACHE.KEY.EMPTY.ERROR")
)
