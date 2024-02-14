package validator

import "fmt"

func SliceRequired[T any](prop string, values []T) Fn {
	return func() error {
		if values == nil {
			return fmt.Errorf("%s must not be nil", prop)
		}
		if len(values) == 0 {
			return fmt.Errorf("%s must not be empty", prop)
		}
		return nil
	}
}

func SliceMaxLength[T any](prop string, values []T, max int) Fn {
	return func() error {
		if len(values) > max {
			return fmt.Errorf("%s is exceeded maximum capacity (max:%d)", prop, max)
		}
		return nil
	}
}

func MapRequired[K comparable, V any](prop string, values map[K]V) Fn {
	return func() error {
		if values == nil {
			return fmt.Errorf("%s must not be nil", prop)
		}
		if len(values) == 0 {
			return fmt.Errorf("%s must not be empty", prop)
		}
		return nil
	}
}
