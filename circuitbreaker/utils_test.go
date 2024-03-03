package circuitbreaker

import (
	"testing"

	"github.com/kanthorlabs/common/clock"
	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/require"
)

func TestDo(t *testing.T) {
	cb, err := NewGoBreaker(testconf, testify.Logger())
	require.NoError(t, err)

	t.Run("OK", func(st *testing.T) {
		expected := testdata.NewUser(clock.New())

		cmd := testdata.Fake.Internet().Domain()
		user, err := Do[testdata.User](cb, cmd, func() (any, error) {
			return &expected, nil
		}, func(err error) error {
			return err
		})

		require.NoError(st, err)
		require.Equal(st, expected, *user)
	})

	t.Run("KO", func(st *testing.T) {
		cmd := testdata.Fake.Internet().Domain()
		_, err := Do[testdata.User](cb, cmd, func() (any, error) {
			return nil, testdata.ErrGeneric
		}, passerror)

		require.ErrorIs(st, err, testdata.ErrGeneric)
	})
}
