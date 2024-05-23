package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func AbsPathify(in string) string {
	if in == "$HOME" || strings.HasPrefix(in, "$HOME"+string(os.PathSeparator)) {
		in = os.Getenv("HOME") + in[5:]
	}

	in = os.ExpandEnv(in)

	if filepath.IsAbs(in) {
		return filepath.Clean(in)
	}

	p, err := filepath.Abs(in)
	if err == nil {
		return filepath.Clean(p)
	}

	return ""
}

func Filepath() (string, error) {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return "", fmt.Errorf("unable to get the file path of the caller")
	}

	return file, nil
}

func Dirpath() (string, error) {
	file, err := Filepath()
	if err != nil {
		return "", err
	}

	return filepath.Dir(file), nil
}
