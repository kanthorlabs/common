package configuration

type Config interface {
	Validate() error
}

type Provider interface {
	Unmarshal(dest any) error
	Sources() []Source
	SetDefault(key string, value any)
	Set(key string, value any)
}

type Source struct {
	Looking string
	Found   string
	Used    bool
}

func New(ns string) (Provider, error) {
	return NewFile(ns, FileLookingDirs)
}
