package signature

import (
	"errors"
	"fmt"
	"strings"
)

var SignatureDivider = ","

var versions map[string]Signature

func init() {
	versions = make(map[string]Signature)
	versions["v1"] = &v1{}
}

func Sign(key, data string) string {
	var signatures []string

	for version := range versions {
		sign := versions[version].Sign(key, data)
		signatures = append(signatures, fmt.Sprintf("%s=%s", version, sign))
	}

	return strings.Join(signatures, SignatureDivider)
}

func Verify(key, data, signature string) error {
	signatures := strings.Split(signature, SignatureDivider)
	for i := range signatures {
		sign := strings.Split(signatures[i], "=")
		if len(sign) != 2 {
			continue
		}

		v, exist := versions[sign[0]]
		if !exist {
			continue
		}

		err := v.Verify(key, data, sign[1])
		if err == nil {
			return nil
		}
	}

	return errors.New("SIGNATURE.VERIFY.NOT_MATCH.ERROR")
}

type Signature interface {
	Sign(key, data string) string
	Verify(key, data, compare string) error
}
