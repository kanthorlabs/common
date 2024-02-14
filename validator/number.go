package validator

import (
	"fmt"
)

func NumberLessThan[T int | int32 | int64 | uint | uint32 | uint64 | float32 | float64](prop string, value T, target T) Fn {
	return func() error {
		if value >= target {
			return fmt.Errorf("%s (%v) must less than %v", prop, value, target)
		}
		return nil
	}
}

func NumberLessThanOrEqual[T int | int32 | int64 | uint | uint32 | uint64 | float32 | float64](prop string, value T, target T) Fn {
	return func() error {
		if value > target {
			return fmt.Errorf("%s (%v) must less than or equal to %v", prop, value, target)
		}
		return nil
	}
}

func NumberGreaterThan[T int | int32 | int64 | uint | uint32 | uint64 | float32 | float64](prop string, value T, target T) Fn {
	return func() error {
		if value <= target {
			return fmt.Errorf("%s (%v) must greater than %v", prop, value, target)
		}
		return nil
	}
}

func NumberGreaterThanOrEqual[T int | int32 | int64 | uint | uint32 | uint64 | float32 | float64](prop string, value T, target T) Fn {
	return func() error {
		if value < target {
			return fmt.Errorf("%s (%v) must greater than or equal to %v", prop, value, target)
		}
		return nil
	}
}

func NumberInRange[T int | int32 | int64 | uint | uint32 | uint64 | float32 | float64](prop string, value T, min, max T) Fn {
	return func() error {
		if value < min {
			return fmt.Errorf("%s (%v) must greater than or equal to %v", prop, value, min)
		}
		if value > max {
			return fmt.Errorf("%s (%v) must less than or equal to %v", prop, value, max)
		}
		return nil
	}
}
