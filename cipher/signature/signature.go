package signature

var SignaturesDivider = ","
var VersionSignatureDivider = "="

type Signature interface {
	Sign(key, data string) string
	Verify(key, data, compare string) error
}

// To prevent downgrade attacks, ignore all schemes that aren’t current support version
// Current support version is v1

var versions map[string]Signature

func init() {
	versions = make(map[string]Signature)
	versions["v1"] = &v1{}
}
