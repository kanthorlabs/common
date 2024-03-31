package webhook

import "errors"

var (
	ErrSignatureMismatch         = errors.New("WEBHOOK.SIGNATURE_MISMATCH.ERROR")
	ErrMessageTimestampMalformed = errors.New("WEBHOOK.MESSAGE.TIMESTAMP_MALFORMED.ERROR")
	ErrMessageTimestampTooOld    = errors.New("WEBHOOK.MESSAGE.TIMESTAMP_TOO_OLD.ERROR")
	ErrMessageTimestampTooNew    = errors.New("WEBHOOK.MESSAGE.TIMESTAMP_TOO_NEW.ERROR")
)
