package entities

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPermission(t *testing.T) {
	t.Run(".Validate", func(st *testing.T) {
		st.Run("KO", func(sst *testing.T) {
			evaluation := &Permission{}
			require.ErrorContains(sst, evaluation.Validate(), "GATEKEEPER.PERMISSION.")
		})
	})
}
