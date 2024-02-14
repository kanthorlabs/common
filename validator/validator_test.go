package validator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	t.Run("stop at first error", func(st *testing.T) {
		err := Validate(StringRequired("name", ""), NumberGreaterThan("age", 0, 0))
		require.ErrorContains(st, err, "is required")
	})
}
