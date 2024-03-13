package validator

import (
	"testing"

	"github.com/kanthorlabs/common/clock"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestPointerNotNil(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		user := testdata.NewUser(clock.New())
		require.NoError(st, PointerNotNil("user", &user)())
	})
	t.Run("KO - nil pointer", func(st *testing.T) {
		require.ErrorContains(st, PointerNotNil[testdata.User]("user", nil)(), "must not be nil")
	})
}

func TestCustom(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		require.NoError(st, Custom("custom", &custom{})())
	})

	t.Run("OK", func(st *testing.T) {
		require.Error(st, Custom("custom", &custom{err: testdata.ErrGeneric})())
	})
}

type custom struct {
	err error
}

func (v *custom) Validate() error {
	return v.err
}
