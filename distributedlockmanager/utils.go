package distributedlockmanager

func Key(k string) (string, error) {
	if k == "" {
		return "", ErrKeyEmpty
	}
	return "dlm/" + k, nil
}
