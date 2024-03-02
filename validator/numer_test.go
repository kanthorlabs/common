package validator

import (
	"testing"

	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

var (
	min = testdata.Fake.IntBetween(256, 512)
	max = min + testdata.Fake.IntBetween(256, 512)
	mid = (min + max) / 2
)

func TestNumberLessThan(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		require.NoError(st, NumberLessThan("bytes", mid, max)())
	})
	t.Run("KO - violate less than rule", func(st *testing.T) {
		require.ErrorContains(st, NumberLessThan("bytes", mid, min)(), "must less than")
	})
}

func TestNumberLessThanOrEqual(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		require.NoError(st, NumberLessThanOrEqual("bytes", mid, max)())
	})
	t.Run("OK - equal", func(st *testing.T) {
		require.NoError(st, NumberLessThanOrEqual("bytes", mid, mid)())
	})
	t.Run("KO - violate less than or equal rule", func(st *testing.T) {
		require.ErrorContains(st, NumberLessThanOrEqual("bytes", mid, min)(), " must less than or equal to")
	})
}

func TestNumberGreaterThan(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		require.NoError(st, NumberGreaterThan("bytes", mid, min)())
	})
	t.Run("KO - violate greater than rule", func(st *testing.T) {
		require.ErrorContains(st, NumberGreaterThan("bytes", mid, max)(), "must greater than")
	})
}

func TestNumberGreaterThanOrEqual(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		require.NoError(st, NumberGreaterThanOrEqual("bytes", mid, min)())
	})
	t.Run("OK - equal", func(st *testing.T) {
		require.NoError(st, NumberGreaterThanOrEqual("bytes", mid, mid)())
	})
	t.Run("KO - violate greater than or equal rule", func(st *testing.T) {
		require.ErrorContains(st, NumberGreaterThanOrEqual("bytes", mid, max)(), " must greater than or equal to")
	})
}

func TestNumberInRange(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		require.NoError(st, NumberInRange("bytes", mid, min, max)())
	})
	t.Run("OK - equal", func(st *testing.T) {
		require.NoError(st, NumberInRange("bytes", mid, mid, max)())
	})
	t.Run("KO - less than", func(st *testing.T) {
		require.ErrorContains(st, NumberInRange("bytes", mid, min/2, min)(), "must less than or equal to")
	})
	t.Run("KO - greater than", func(st *testing.T) {
		require.ErrorContains(st, NumberInRange("bytes", mid, max, max*2)(), "must greater than or equal to")
	})
}
