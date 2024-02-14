package validator

import (
	"testing"

	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/timer"
	"github.com/stretchr/testify/require"
)

func TestPointerNotNil(t *testing.T) {
	t.Run("ok", func(st *testing.T) {
		require.Nil(st, PointerNotNil("user", testdata.NewUser(timer.New()))())
	})
	t.Run("ko", func(st *testing.T) {
		require.ErrorContains(st, PointerNotNil[testdata.User]("user", nil)(), "must not be nil")
	})
}
