package safe

import (
	"database/sql/driver"
	"encoding/json"
	"sync"
)

type Metadata struct {
	kv map[string]any
	mu sync.Mutex
}

func (meta *Metadata) Set(k string, v any) {
	meta.mu.Lock()
	defer meta.mu.Unlock()

	if meta.kv == nil {
		meta.kv = make(map[string]any)
	}

	meta.kv[k] = v
}

func (meta *Metadata) Get(k string) (any, bool) {
	meta.mu.Lock()
	defer meta.mu.Unlock()

	if meta.kv == nil {
		return nil, false
	}

	v, has := meta.kv[k]
	return v, has
}

func (meta *Metadata) Merge(src *Metadata) {
	meta.mu.Lock()
	defer meta.mu.Unlock()

	if meta.kv == nil {
		meta.kv = make(map[string]any)
	}

	if len(src.kv) == 0 {
		return
	}

	for k := range src.kv {
		meta.kv[k] = src.kv[k]
	}
}

func (meta *Metadata) String() string {
	data, _ := json.Marshal(meta.kv)
	return string(data)
}

// Value implements the driver Valuer interface.
func (meta *Metadata) Value() (driver.Value, error) {
	data, err := json.Marshal(meta.kv)
	return string(data), err
}

// Scan implements the Scanner interface.
func (meta *Metadata) Scan(value any) error {
	return json.Unmarshal([]byte(value.(string)), &meta.kv)
}
