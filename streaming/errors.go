package streaming

import "errors"

var (
	ErrNotConnected        = errors.New("STREAMING.NOT_CONNECTED.ERROR")
	ErrAlreadyConnected    = errors.New("STREAMING.ALREADY_CONNECTED.ERROR")
	ErrSubNotConnected     = errors.New("STREAMING.SUBSCRIBER.NOT_CONNECTED.ERROR")
	ErrSubAlreadyConnected = errors.New("STREAMING.SUBSCRIBER.ALREADY_CONNECTED.ERROR")
	ErrSubTerminiated      = errors.New("STREAMING.SUBSCRIBER.TERMINATED.ERROR")
	ErrSubAckFail          = errors.New("STREAMING.SUBSCRIBER.ACK_FAIL.ERROR")
	ErrSubNakFail          = errors.New("STREAMING.SUBSCRIBER.NAK_FAIL.ERROR")
)
