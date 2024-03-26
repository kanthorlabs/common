package strategies

import (
	"encoding/base64"
	"strings"

	"github.com/kanthorlabs/common/passport/entities"
)

var (
	SchemeBasic   = "Basic "
	RegionDivider = "."
)

func IsBasicScheme(raw string) bool {
	if len(raw) < len(SchemeBasic) {
		return false
	}
	if !strings.EqualFold(raw[:len(SchemeBasic)], SchemeBasic) {
		return false
	}
	return true
}

func ParseBasicCredentials(raw string) (*entities.Credentials, error) {
	if !IsBasicScheme(raw) {
		return nil, ErrParseCredentials
	}

	c, err := base64.StdEncoding.DecodeString(raw[len(SchemeBasic):])
	if err != nil {
		return nil, ErrParseCredentials
	}
	cs := string(c)

	if i := strings.Index(cs, RegionDivider); i >= 0 {
		region := cs[:i]
		credentials := cs[i+len(RegionDivider):]

		username, password, ok := strings.Cut(credentials, ":")
		if !ok {
			return nil, ErrParseCredentials
		}

		return &entities.Credentials{
			Username: username,
			Password: password,
			Region:   region,
		}, nil
	}

	username, password, ok := strings.Cut(cs, ":")
	if !ok {
		return nil, ErrParseCredentials
	}

	return &entities.Credentials{Username: username, Password: password}, nil
}
