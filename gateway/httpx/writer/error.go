package writer

func Error(err error) M {
	return M{"error": err.Error()}
}

func ErrorString(message string) M {
	return M{"error": message}
}
