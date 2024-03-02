package signature

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

// v1 is the implementation of the signature version 1 that is based on HMAC-SHA256 with hex encoding.
type v1 struct{}

func (signature *v1) Sign(key, data string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))

	return hex.EncodeToString(mac.Sum(nil))
}

func (signature *v1) Verify(key, data, compare string) error {
	sign := signature.Sign(key, data)

	if hmac.Equal([]byte(sign), []byte(compare)) {
		return nil
	}

	return errors.New("SIGNATURE.V1.VERIFY.NOT_MATCH.ERROR")
}
