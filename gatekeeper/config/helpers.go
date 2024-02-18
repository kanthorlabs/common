package config

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/kanthorlabs/common/gatekeeper/entities"
)

func ParseDefinitionsToPermissions(uri string) (map[string][]entities.Permission, error) {
	var definitions map[string][]entities.Permission

	if strings.HasPrefix(uri, "file://") {
		p := strings.Replace(uri, "file://", "", -1)
		data, err := os.ReadFile(p)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(data, &definitions); err != nil {
			return nil, err
		}

		return definitions, nil
	}

	if strings.HasPrefix(uri, "base64://") {
		bs := strings.Replace(uri, "base64://", "", -1)
		data, err := base64.StdEncoding.DecodeString(bs)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(data, &definitions); err != nil {
			return nil, err
		}

		return definitions, nil
	}

	return nil, errors.New("GATEKEEPER.CONFIG.DEFINITIONS.URI.UNSUPPORTED.ERROR")
}
