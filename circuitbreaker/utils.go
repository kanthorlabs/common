package circuitbreaker

// Do is a helper function that wraps the CircuitBreaker.Do method and returns the result as the expected type.
func Do[T any](cb CircuitBreaker, cmd string, onHandle Handler, onError ErrorHandler) (*T, error) {
	out, err := cb.Do(cmd, onHandle, onError)

	if err != nil {
		return nil, err
	}

	return out.(*T), nil
}
