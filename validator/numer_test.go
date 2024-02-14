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
	t.Run("ok", func(st *testing.T) {
		require.Nil(st, NumberLessThan("bytes", mid, max)())
	})
	t.Run("ko", func(st *testing.T) {
		require.ErrorContains(st, NumberLessThan("bytes", mid, min)(), "must less than")
	})
}

func TestNumberLessThanOrEqual(t *testing.T) {
	t.Run("ok", func(st *testing.T) {
		require.Nil(st, NumberLessThanOrEqual("bytes", mid, max)())
	})
	t.Run("ok because of equal", func(st *testing.T) {
		require.Nil(st, NumberLessThanOrEqual("bytes", mid, mid)())
	})
	t.Run("ko", func(st *testing.T) {
		require.ErrorContains(st, NumberLessThanOrEqual("bytes", mid, min)(), " must less than or equal to")
	})
}

func TestNumberGreaterThan(t *testing.T) {
	t.Run("ok", func(st *testing.T) {
		require.Nil(st, NumberGreaterThan("bytes", mid, min)())
	})
	t.Run("ko", func(st *testing.T) {
		require.ErrorContains(st, NumberGreaterThan("bytes", mid, max)(), "must greater than")
	})
}

func TestNumberGreaterThanOrEqual(t *testing.T) {
	t.Run("ok", func(st *testing.T) {
		require.Nil(st, NumberGreaterThanOrEqual("bytes", mid, min)())
	})
	t.Run("ok because of equal", func(st *testing.T) {
		require.Nil(st, NumberGreaterThanOrEqual("bytes", mid, mid)())
	})
	t.Run("ko", func(st *testing.T) {
		require.ErrorContains(st, NumberGreaterThanOrEqual("bytes", mid, max)(), " must greater than or equal to")
	})
}

func TestNumberInRange(t *testing.T) {
	t.Run("ok", func(st *testing.T) {
		require.Nil(st, NumberInRange("bytes", mid, min, max)())
	})
	t.Run("ok because of equal", func(st *testing.T) {
		require.Nil(st, NumberInRange("bytes", mid, mid, max)())
	})
	t.Run("ko because of greater than", func(st *testing.T) {
		require.ErrorContains(st, NumberInRange("bytes", mid, max, max*2)(), "must greater than or equal to")
	})
	t.Run("ko because of less than", func(st *testing.T) {
		require.ErrorContains(st, NumberInRange("bytes", mid, min/2, min)(), "must less than or equal to")
	})
}
