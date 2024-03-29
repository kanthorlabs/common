package utils

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/kanthorlabs/common/passport/entities"
	"github.com/kanthorlabs/common/project"
)

var (
	SchemeBasic   = "Basic "
	RegionDivider = " "
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

func CreateRegionalBasicCredentials(raw string) string {
	return base64.StdEncoding.EncodeToString([]byte(raw + RegionDivider + project.Region()))
}

func ParseBasicCredentials(raw string) (*entities.Credentials, error) {
	if !IsBasicScheme(raw) {
		return nil, errors.New("PASSPORT.UTILS.PARSE_BASIC_CREDENTIALS.SCHEME_UNKNOWN.ERROR")
	}

	c, err := base64.StdEncoding.DecodeString(raw[len(SchemeBasic):])
	if err != nil {
		return nil, errors.New("PASSPORT.UTILS.PARSE_BASIC_CREDENTIALS.DECODE.ERROR")
	}
	cs := string(c)

	credentials, region, _ := strings.Cut(cs, RegionDivider)
	username, password, ok := strings.Cut(credentials, ":")
	if !ok {
		return nil, errors.New("PASSPORT.UTILS.PARSE_BASIC_CREDENTIALS.PARSE.ERROR")
	}

	return &entities.Credentials{
		Username: username,
		Password: password,
		Region:   region,
	}, nil
}
