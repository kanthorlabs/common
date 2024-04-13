package configuration

import (
	"testing"

	"github.com/kanthorlabs/common/project"
	"github.com/kanthorlabs/common/safe"
	"github.com/stretchr/testify/require"
)

type configs struct {
	Counter  int            `json:"counter" yaml:"counter" mapstructure:"counter"`
	Blood    string         `json:"blood,omitempty" yaml:"blood,omitempty" mapstructure:"blood"`
	Metadata *safe.Metadata `json:"metadata,omitempty" yaml:"metadata,omitempty" mapstructure:"metadata"`

	UnsetByDefault string `json:"unset_by_default,omitempty" yaml:"unset_by_default,omitempty" mapstructure:"unset_by_default"`
}

func TestNew(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		provider, err := New(project.Namespace())
		require.NoError(st, err)
		require.NotNil(st, provider)
	})
}
