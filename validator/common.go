package validator

import (
	"fmt"
)

func PointerNotNil[T any](prop string, value *T) Fn {
	return func() error {
		if value == nil {
			return fmt.Errorf("%s must not be nil", prop)
		}
		return nil
	}
}
