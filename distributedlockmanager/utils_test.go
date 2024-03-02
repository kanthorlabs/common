package distributedlockmanager

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKey(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		k, err := Key("test")
		require.NoError(st, err)
		require.Equal(st, "dlm/test", k)
	})

	t.Run("KO - key empty error", func(st *testing.T) {
		_, err := Key("")
		require.ErrorIs(st, err, ErrKeyEmpty)
	})
}
