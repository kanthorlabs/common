package validator

func Slice[T any](items []T, fn func(i int, item *T) error) Fn {
	return func() error {
		for i := 0; i < len(items); i++ {
			if err := fn(i, &items[i]); err != nil {
				return err
			}
		}

		return nil
	}
}

func Map[T any](items map[string]T, fn func(refId string, item T) error) Fn {
	return func() error {
		for refId := range items {
			if err := fn(refId, items[refId]); err != nil {
				return err
			}
		}

		return nil
	}
}
