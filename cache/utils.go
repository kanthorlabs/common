package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

func GetOrSet[T any](cache Cache, ctx context.Context, key string, ttl time.Duration, fn func() (*T, error)) (*T, error) {
	entry, err := cache.Get(ctx, key)
	if err == nil {
		var dest T
		if err := json.Unmarshal(entry, &dest); err != nil {
			return nil, err
		}

		return &dest, nil
	}

	// if we catched any error other than ErrEntryNotFound, return it immediately
	if !errors.Is(err, ErrEntryNotFound) {
		return nil, err
	}

	data, err := fn()
	if err != nil {
		return nil, err
	}

	entry, err = json.Marshal(data)
	if err != nil {
		return nil, err
	}

	if err := cache.Set(ctx, key, entry, ttl); err != nil {
		return nil, err
	}

	return data, nil
}
