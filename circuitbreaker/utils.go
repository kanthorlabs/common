package circuitbreaker

func Do[T any](cb CircuitBreaker, cmd string, onHandle Handler, onError ErrorHandler) (*T, error) {
	out, err := cb.Do(cmd, onHandle, onError)

	if out == nil {
		return nil, err
	}

	return out.(*T), nil
}
