package streaming

import "errors"

var (
	ErrNotConnected        = errors.New("STREAMING.NOT_CONNECTED.ERROR")
	ErrAlreadyConnected    = errors.New("STREAMING.ALREADY_CONNECTED.ERROR")
	ErrPubEventValidation  = errors.New("STREAMING.PUBLISHER.EVENT_VALIDATION.ERROR")
	ErrPubEventPublish     = errors.New("STREAMING.PUBLISHER.EVENT_PUBLISH.ERROR")
	ErrPubEventDuplicated  = errors.New("STREAMING.PUBLISHER.EVENT_DUPLICATED.ERROR")
	ErrSubNotConnected     = errors.New("STREAMING.SUBSCRIBER.NOT_CONNECTED.ERROR")
	ErrSubAlreadyConnected = errors.New("STREAMING.SUBSCRIBER.ALREADY_CONNECTED.ERROR")
	ErrSubAckFail          = errors.New("STREAMING.SUBSCRIBER.ACK_FAIL.ERROR")
	ErrSubNakFail          = errors.New("STREAMING.SUBSCRIBER.NAK_FAIL.ERROR")
)
