package signature

import (
	"errors"
	"strings"
)

// Verify verifies the data with the provided key using all available versions.
// It compares the provided signature with all available versions.
// If any of the versions match, it returns no error, otherwise it returns an error of type "SIGNATURE.VERIFY.NOT_MATCH.ERROR".
func Verify(key, data, signature string) error {
	signatures := strings.Split(signature, SignaturesDivider)
	for i := range signatures {
		versionsign := strings.Split(signatures[i], VersionSignatureDivider)
		if len(versionsign) != 2 {
			continue
		}

		v, exist := versions[versionsign[0]]
		if !exist {
			continue
		}

		err := v.Verify(key, data, versionsign[1])
		if err == nil {
			return nil
		}
	}

	return errors.New("SIGNATURE.VERIFY.NOT_MATCH.ERROR")
}

// VerifyAny verifies the signature with the provided keys using all available versions.
// The usecase is you want to rotate the key and still be able to verify the old signature
func VerifyAny(keys []string, data, signature string) error {
	for _, key := range keys {
		err := Verify(key, data, signature)
		if err == nil {
			return nil
		}
	}
	return errors.New("SIGNATURE.VERIFY.NOT_MATCH.ERROR")
}
