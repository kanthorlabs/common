package validator

import (
	"testing"

	"github.com/kanthorlabs/common/clock"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestPointerNotNil(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		require.Nil(st, PointerNotNil("user", testdata.NewUser(clock.New()))())
	})
	t.Run("KO - nil pointer", func(st *testing.T) {
		require.ErrorContains(st, PointerNotNil[testdata.User]("user", nil)(), "must not be nil")
	})
}
