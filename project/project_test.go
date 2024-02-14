package project

import (
	"fmt"
	"testing"

	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestVersion(t *testing.T) {
	v := "v2024.531.530"
	SetVersion("v2024.531.530")
	require.Equal(t, v, GetVersion())
}

func TestEnv(t *testing.T) {
	require.False(t, IsDev())
	t.Setenv("PROJECT_ENV", DevEnv)
	require.True(t, IsDev())
}

func TestName(t *testing.T) {
	name := testdata.Fake.Lorem().Word()

	t.Run("default", func(st *testing.T) {
		got := Name(name)
		require.Equal(st, fmt.Sprintf("%s_%s", Namespace(), name), got)
	})
}

func TestTopic(t *testing.T) {
	first := testdata.Fake.Lorem().Word()
	second := testdata.Fake.Lorem().Word()

	t.Run("default", func(st *testing.T) {
		got := Topic(first, second)
		require.Equal(st, fmt.Sprintf("%s.%s", first, second), got)
	})

	t.Run("with empty string", func(st *testing.T) {
		got := Topic(first, "", second)
		require.Equal(st, fmt.Sprintf("%s.%s", first, second), got)
	})
}

func TestIsTopic(t *testing.T) {
	first := testdata.Fake.Lorem().Word()
	second := testdata.Fake.Lorem().Word()

	t.Run("true", func(st *testing.T) {
		got := IsTopic(Subject(first, second), first)
		require.True(st, got)
	})

	t.Run("false", func(st *testing.T) {
		got := IsTopic(Subject(first, second), second)
		require.False(st, got)
	})
}

func TestSubject(t *testing.T) {
	first := testdata.Fake.Lorem().Word()
	second := testdata.Fake.Lorem().Word()

	t.Run("default", func(st *testing.T) {
		got := Subject(first, second)
		require.Equal(st, fmt.Sprintf("%s.%s.%s.%s.%s", Namespace(), Region(), Tier(), first, second), got)
	})

	t.Run("with empty string", func(st *testing.T) {
		got := Subject(first, "", second)
		require.Equal(st, fmt.Sprintf("%s.%s.%s.%s.%s", Namespace(), Region(), Tier(), first, second), got)
	})
}
