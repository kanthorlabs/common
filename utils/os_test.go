package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAbsPathify(t *testing.T) {
	home := os.Getenv("HOME")
	homer := filepath.Join(home, "homer")
	wd, _ := os.Getwd()

	t.Setenv("HOMER_ABSOLUTE_PATH", homer)
	t.Setenv("VAR_WITH_RELATIVE_PATH", "relative")

	tests := []struct {
		input  string
		output string
	}{
		{"", wd},
		{"sub", filepath.Join(wd, "sub")},
		{"./", wd},
		{"./sub", filepath.Join(wd, "sub")},
		{"$HOME", home},
		{"$HOME/", home},
		{"$HOME/sub", filepath.Join(home, "sub")},
		{"$HOMER_ABSOLUTE_PATH", homer},
		{"$HOMER_ABSOLUTE_PATH/", homer},
		{"$HOMER_ABSOLUTE_PATH/sub", filepath.Join(homer, "sub")},
		{"$VAR_WITH_RELATIVE_PATH", filepath.Join(wd, "relative")},
		{"$VAR_WITH_RELATIVE_PATH/", filepath.Join(wd, "relative")},
		{"$VAR_WITH_RELATIVE_PATH/sub", filepath.Join(wd, "relative", "sub")},
	}

	for _, test := range tests {
		got := AbsPathify(test.input)
		require.Equal(t, test.output, got)
	}
}

func TestFilepath(t *testing.T) {
	file, err := Filepath()
	require.NoError(t, err)
	require.NotEmpty(t, file)
	require.True(t, strings.HasSuffix(file, "/utils/os_test.go"))
}

func TestDirpath(t *testing.T) {
	file, err := Dirpath()
	require.NoError(t, err)
	require.NotEmpty(t, file)
	require.True(t, strings.HasSuffix(file, "/utils"))
}
