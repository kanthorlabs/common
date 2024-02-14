package configuration

type Config interface {
	Validate() error
}

type Provider interface {
	Unmarshal(dest interface{}) error
	Sources() []Source
	SetDefault(key string, value interface{})
	Set(key string, value interface{})
}

type Source struct {
	Looking string
	Found   string
	Used    bool
}

func New(ns string) (Provider, error) {
	return NewFile(ns, FileLookingDirs)
}
