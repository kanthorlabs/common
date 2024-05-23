package cache

import (
	"encoding/json"
	"fmt"
)

func Marshal(v any) ([]byte, error) {
	if v == nil {
		return []byte{}, nil
	}

	var entry []byte
	entry, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("CACHE.VALUE.MARSHAL.ERROR: %w", err)
	}

	return entry, nil
}

func Unmarshal(data []byte, v any) error {
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("CACHE.VALUE.UNMARSHAL.ERROR: %w", err)
	}

	return nil
}
