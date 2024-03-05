package writer

func Error(err error) E {
	return E{Error: err.Error()}
}

func ErrorString(message string) E {
	return E{Error: message}
}
