package distributedlockmanager

import (
	"testing"

	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/require"
)

func TestKey(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		k, err := Key("test")
		require.NoError(st, err)
		require.Equal(st, "dlm/test", k)
	})

	t.Run(testify.CaseKOKeyEmptyError, func(st *testing.T) {
		_, err := Key("")
		require.ErrorIs(st, err, ErrKeyEmpty)
	})
}
