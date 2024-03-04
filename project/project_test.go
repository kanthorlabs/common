package project

import (
	"fmt"
	"strings"
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
	t.Setenv("KANTHOR_ENV", DevEnv)
	require.True(t, IsDev())
}

func TestName(t *testing.T) {
	name := testdata.Fake.Lorem().Word()
	require.Equal(t, fmt.Sprintf("%s_%s", Namespace(), name), Name(name))
}

func TestTopic(t *testing.T) {
	first := testdata.Fake.Lorem().Word()
	second := testdata.Fake.Lorem().Word()

	t.Run("default", func(st *testing.T) {
		require.Equal(st, fmt.Sprintf("%s.%s", first, second), Topic(first, second))
	})

	t.Run("with empty string", func(st *testing.T) {
		require.Equal(st, fmt.Sprintf("%s.%s", first, second), Topic(first, second))
	})
}

func TestIsTopic(t *testing.T) {
	first := testdata.Fake.Lorem().Word()
	second := strings.Join(testdata.Fake.Lorem().Words(10), " ")

	t.Run("true", func(st *testing.T) {
		require.True(st, IsTopic(Subject(first, second), first))
	})

	t.Run("false", func(st *testing.T) {
		require.False(st, IsTopic(Subject(first, second), second))
	})
}

func TestSubject(t *testing.T) {
	first := testdata.Fake.Lorem().Word()
	second := testdata.Fake.Lorem().Word()

	t.Run("default", func(st *testing.T) {
		subject := Subject(first, second)
		require.Equal(st, fmt.Sprintf("%s.%s.%s.%s.%s", Namespace(), Region(), Tier(), first, second), subject)
	})

	t.Run("with empty string", func(st *testing.T) {
		subject := Subject(first, "", second)
		require.Equal(st, fmt.Sprintf("%s.%s.%s.%s.%s", Namespace(), Region(), Tier(), first, second), subject)
	})
}
