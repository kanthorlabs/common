package idempotency

func Key(k string) (string, error) {
	if k == "" {
		return "", ErrKeyEmpty
	}
	return "idempotency/" + k, nil
}
