package signature

import "strings"

// Sign signs the data with the provided key using all available versions.
// The result is a string with all signatures with format "version=signature" separated by commas.
func Sign(key, data string) string {
	var signatures []string

	for version := range versions {
		sign := versions[version].Sign(key, data)
		signatures = append(signatures, version+VersionSignatureDivider+sign)
	}

	return strings.Join(signatures, SignaturesDivider)
}
