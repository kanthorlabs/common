package distributedlockmanager

func Key(k string) string {
	return "dlm/" + k
}
