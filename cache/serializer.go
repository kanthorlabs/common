package cache

import (
	"encoding/json"
	"fmt"
)

// Marshal is a helper function to marshal a value into a byte slice.
// Use json.Marshal under the hood.
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

// Unmarshal is a helper function to unmarshal a byte slice into a value.
// Use json.Unmarshal under the hood.
func Unmarshal(data []byte, v any) error {
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("CACHE.VALUE.UNMARSHAL.ERROR: %w", err)
	}

	return nil
}
