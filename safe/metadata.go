package safe

import (
	"database/sql/driver"
	"encoding/json"
	"sync"

	"gopkg.in/yaml.v3"
)

// Metadata is a thread-safe key-value store.
// There are some gotchas when using this struct:
// json marshalling is supported but all number will be parsed as float64.
// yaml marshalling is supported but all number will be parsed as int.
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

	if src == nil || len(src.kv) == 0 {
		return
	}

	for k := range src.kv {
		meta.kv[k] = src.kv[k]
	}
}

func (meta *Metadata) String() string {
	if meta == nil || meta.kv == nil {
		return ""
	}

	data, _ := json.Marshal(meta.kv)
	return string(data)
}

// Value implements the driver Valuer interface.
func (meta *Metadata) Value() (driver.Value, error) {
	// meta == nil when we convert it to database value
	if meta == nil || meta.kv == nil {
		return "", nil
	}
	data, err := json.Marshal(meta.kv)
	return string(data), err
}

// Scan implements the Scanner interface.
func (meta *Metadata) Scan(value any) error {
	v := value.(string)
	if v == "" {
		return nil
	}
	return json.Unmarshal([]byte(v), &meta.kv)
}

func (meta *Metadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(meta.kv)
}

func (meta *Metadata) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &meta.kv)
}

func (meta *Metadata) MarshalYAML() (interface{}, error) {
	var value yaml.Node
	return value, value.Encode(meta.kv)
}

func (meta *Metadata) UnmarshalYAML(value *yaml.Node) error {
	return value.Decode(&meta.kv)
}
