package configuration

import (
	"os"
	"path"
	"testing"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/project"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

type configs struct {
	Counter int
	Blood   string
}

func TestFile(t *testing.T) {
	confName := "blood"

	t.Run("$KANTHOR_HOME", func(st *testing.T) {
		wd, err := os.Getwd()
		require.Nil(st, err)

		st.Setenv("KANTHOR_HOME", path.Join(wd, "..", ".kanthor"))
		provider, err := New(project.Namespace())
		require.Nil(st, err)

		var conf configs
		require.Nil(st, provider.Unmarshal(&conf))
		require.True(st, provider.Sources()[0].Used)
	})

	t.Run("$HOME/.kanthor", func(st *testing.T) {
		wd, err := os.Getwd()
		require.Nil(st, err)

		st.Setenv("HOME", path.Join(wd, ".."))
		provider, err := New(project.Namespace())
		require.Nil(st, err)

		var conf configs
		require.Nil(st, provider.Unmarshal(&conf))

		require.Equal(st, 1, conf.Counter)
		require.True(st, provider.Sources()[1].Used)
	})

	t.Run("current working directory", func(st *testing.T) {
		provider, err := New(project.Namespace())
		require.Nil(st, err)

		var conf configs
		require.Nil(st, provider.Unmarshal(&conf))

		require.Equal(st, 2, conf.Counter)
		require.True(st, provider.Sources()[2].Used)
	})

	t.Run(".SetDefault", func(st *testing.T) {
		provider, err := New(project.Namespace())
		require.Nil(st, err)

		blood := testdata.Fake.Blood().Name()
		provider.SetDefault(confName, blood)

		var conf configs
		require.Nil(st, provider.Unmarshal(&conf))

		require.Equal(st, blood, conf.Blood)
	})

	t.Run(".Set must override .SetDefault", func(st *testing.T) {
		provider, err := New(project.Namespace())
		require.Nil(st, err)

		blood := testdata.Fake.Blood().Name()
		provider.SetDefault(confName, blood)

		override := uuid.NewString()
		provider.Set(confName, override)
		var conf configs
		require.Nil(st, provider.Unmarshal(&conf))

		require.Equal(st, override, conf.Blood)
	})
}
