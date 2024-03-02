package cache

import (
	"context"
	"errors"
	"time"
)

func Key(k string) (string, error) {
	if k == "" {
		return "", ErrKeyEmpty
	}
	return "cache/" + k, nil
}

// GetOrSet is a helper function that allow you get existing entry from the cache or set it if it does not exist yet
// the function must return a pointer of the expected type and an error
func GetOrSet[T any](cache Cache, ctx context.Context, key string, ttl time.Duration, fn func() (*T, error)) (*T, error) {
	var err error
	var dest T

	err = cache.Get(ctx, key, &dest)
	if err == nil {
		return &dest, nil
	}

	// if the error is returned and it is not ErrEntryNotFound, return the error
	if !errors.Is(err, ErrEntryNotFound) {
		return nil, err
	}

	// otherwise, retrieve the entry and set it in the cache
	entry, err := fn()
	if err != nil {
		return nil, err
	}

	return entry, cache.Set(ctx, key, entry, ttl)
}
