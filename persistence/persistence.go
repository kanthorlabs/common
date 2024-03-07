package persistence

import "github.com/kanthorlabs/common/patterns"

type Persistence interface {
	patterns.Connectable
	Client() any
}
