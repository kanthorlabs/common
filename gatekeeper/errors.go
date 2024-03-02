package gatekeeper

import "errors"

var (
	ErrNotReady         = errors.New("GATEKEEPER.NOT_READY.ERROR")
	ErrNotLive          = errors.New("GATEKEEPER.NOT_LIVE.ERROR")
	ErrAlreadyConnected = errors.New("GATEKEEPER.ALREADY_CONNECTED.ERROR")
	ErrNotConnected     = errors.New("GATEKEEPER.NOT_CONNECTED.ERROR")
	ErrPrivilegeQuery   = errors.New("GATEKEEPER.PRIVILEGE.QUERY.ERROR")
)
