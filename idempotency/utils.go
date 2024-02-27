package idempotency

func Key(k string) string {
	return "idempotency/" + k
}
