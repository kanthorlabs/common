package utils

import (
	"runtime"
	"testing"
)

// Skip some tests on Windows that kept failing when Windows was added to the CI as a target.
//
//nolint:gocritic // sloppyTestFuncName
func skipWindows(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skip test on Windows")
	}
}
