package entities

import (
	"testing"

	"github.com/kanthorlabs/common/project"
	"github.com/stretchr/testify/require"
)

func TestPrivilege(t *testing.T) {
	t.Run(".TableName", func(st *testing.T) {
		st.Run("KO", func(sst *testing.T) {
			evaluation := &Privilege{}
			require.Equal(sst, evaluation.TableName(), project.Name("gatekeeper_privilege"))
		})
	})

	t.Run(".Validate", func(st *testing.T) {
		st.Run("KO", func(sst *testing.T) {
			evaluation := &Privilege{}
			require.ErrorContains(sst, evaluation.Validate(), "GATEKEEPER.PREVILEGE.")
		})
	})
}
