package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKey(t *testing.T) {
	tests := map[string][]string{
		"":      {""},
		"a":     {"a"},
		"a/b/c": {"a", "b", "c"},
	}

	for expected, values := range tests {
		got := Key(values...)
		require.Equal(t, expected, got)
	}
}

func TestStringify(t *testing.T) {
	tests := map[string]any{
		"":   nil,
		"{}": map[string]any{},
	}

	for expected, value := range tests {
		got := Stringify(value)
		require.Equal(t, expected, got)
	}
}

func TestStringifyIndent(t *testing.T) {
	tests := map[string]any{
		"":                     nil,
		"{}":                   map[string]any{},
		"{\n  \"ok\": true\n}": map[string]any{"ok": true},
	}

	for expected, value := range tests {
		got := StringifyIndent(value)
		require.Equal(t, expected, got)
	}
}

func TestRandomString(t *testing.T) {
	t.Run("31 chars", func(st *testing.T) {
		got := RandomString(31)
		require.Equal(st, len(got), 31)
	})

	t.Run("32 chars", func(st *testing.T) {
		got := RandomString(32)
		require.Equal(st, len(got), 32)
	})

	t.Run("65 chars", func(st *testing.T) {
		got := RandomString(65)
		require.Equal(st, len(got), 65)
	})
}
