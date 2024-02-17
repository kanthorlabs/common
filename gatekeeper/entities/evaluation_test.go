package entities

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEvaluation(t *testing.T) {
	t.Run(".Validate", func(st *testing.T) {
		st.Run("KO", func(sst *testing.T) {
			evaluation := &Evaluation{}
			require.ErrorContains(sst, evaluation.Validate(), "GATEKEEPER.EVALUATION.")
		})
	})
}
