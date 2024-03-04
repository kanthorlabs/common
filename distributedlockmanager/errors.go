package distributedlockmanager

import "errors"

var (
	ErrAlreadyConnected = errors.New("DISTRIBUTED_LOCK_MANAGER.ALREADY_CONNECTED.ERROR")
	ErrNotConnected     = errors.New("DISTRIBUTED_LOCK_MANAGER.NOT_CONNECTED.ERROR")
	ErrKeyEmpty         = errors.New("DISTRIBUTED_LOCK_MANAGER.KEY.EMPTY.ERROR")
	ErrLock             = errors.New("DISTRIBUTED_LOCK_MANAGER.LOCK.ERROR")
	ErrUnlock           = errors.New("DISTRIBUTED_LOCK_MANAGER.UNLOCK.ERROR")
)
