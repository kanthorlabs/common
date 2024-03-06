package configuration

import (
	"testing"

	"github.com/kanthorlabs/common/project"
	"github.com/kanthorlabs/common/safe"
	"github.com/stretchr/testify/require"
)

type configs struct {
	Counter        int            `json:"counter" yaml:"counter" mapstructure:"counter"`
	Blood          string         `json:"blood,omitempty" yaml:"blood,omitempty" mapstructure:"blood"`
	UnsetByDefault string         `json:"unset_by_default,omitempty" yaml:"unset_by_default,omitempty" mapstructure:"unset_by_default"`
	Metadata       *safe.Metadata `json:"metadata,omitempty" yaml:"metadata,omitempty" mapstructure:"metadata"`
}

func TestNew(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		provider, err := New(project.Namespace())
		require.NoError(st, err)
		require.NotNil(st, provider)
	})
}

// func TestFile(t *testing.T) {
// 	confName := "blood"

// 	t.Run("OK - $KANTHOR_HOME", func(st *testing.T) {
// 		wd, err := os.Getwd()
// 		require.NoError(st, err)

// 		st.Setenv("KANTHOR_HOME", path.Join(wd, "..", ".kanthor"))
// 		provider, err := New(project.Namespace())
// 		require.NoError(st, err)

// 		var conf configs
// 		require.NoError(st, provider.Unmarshal(&conf))
// 		require.True(st, provider.Sources()[0].Used)
// 	})

// 	t.Run("OK - $HOME/.kanthor", func(st *testing.T) {
// 		wd, err := os.Getwd()
// 		require.NoError(st, err)

// 		st.Setenv("HOME", path.Join(wd, ".."))
// 		provider, err := New(project.Namespace())
// 		require.NoError(st, err)

// 		var conf configs
// 		require.NoError(st, provider.Unmarshal(&conf))

// 		require.Equal(st, 1, conf.Counter)
// 		require.True(st, provider.Sources()[1].Used)
// 	})

// 	t.Run("OK - current working directory", func(st *testing.T) {
// 		provider, err := New(project.Namespace())
// 		require.NoError(st, err)

// 		var conf configs
// 		require.NoError(st, provider.Unmarshal(&conf))

// 		require.Equal(st, 2, conf.Counter)
// 		require.True(st, provider.Sources()[2].Used)
// 	})

// 	t.Run(".SetDefault", func(st *testing.T) {
// 		provider, err := New(project.Namespace())
// 		require.NoError(st, err)

// 		blood := testdata.Fake.Blood().Name()
// 		provider.SetDefault(confName, blood)

// 		var conf configs
// 		require.NoError(st, provider.Unmarshal(&conf))

// 		require.Equal(st, blood, conf.Blood)
// 	})

// 	t.Run(".Set must override .SetDefault", func(st *testing.T) {
// 		provider, err := New(project.Namespace())
// 		require.NoError(st, err)

// 		blood := testdata.Fake.Blood().Name()
// 		provider.SetDefault(confName, blood)

// 		override := uuid.NewString()
// 		provider.Set(confName, override)
// 		var conf configs
// 		require.NoError(st, provider.Unmarshal(&conf))

// 		require.Equal(st, override, conf.Blood)
// 	})
// }
