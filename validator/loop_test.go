package validator

import (
	"testing"

	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestSliceRequired(t *testing.T) {
	t.Run("ok", func(st *testing.T) {
		agents := []string{
			testdata.Fake.UserAgent().Chrome(),
			testdata.Fake.UserAgent().Firefox(),
			testdata.Fake.UserAgent().Safari(),
		}
		require.Nil(st, SliceRequired("agents", agents)())
	})

	t.Run("ko because of nil", func(st *testing.T) {
		require.ErrorContains(st, SliceRequired[string]("agents", nil)(), "must not be nil")
	})

	t.Run("ko because of empty", func(st *testing.T) {
		require.ErrorContains(st, SliceRequired[string]("agents", []string{})(), "must not be empty")
	})
}

func TestSliceMaxLength(t *testing.T) {
	count := testdata.Fake.IntBetween(100000, 999999)
	items := make([]int, count)

	t.Run("ok", func(st *testing.T) {
		require.Nil(st, SliceMaxLength("items", items, count)())
	})

	t.Run("ko", func(st *testing.T) {
		require.ErrorContains(st, SliceMaxLength("items", items, count-1)(), "is exceeded maximum capacity")
	})
}

func TestMapRequired(t *testing.T) {
	t.Run("ok", func(st *testing.T) {
		agents := map[string]string{
			"chrome":  testdata.Fake.UserAgent().Chrome(),
			"firefox": testdata.Fake.UserAgent().Firefox(),
			"safar":   testdata.Fake.UserAgent().Safari(),
		}
		require.Nil(st, MapRequired("agents", agents)())
	})

	t.Run("ko because of nil", func(st *testing.T) {
		require.ErrorContains(st, MapRequired[string, string]("agents", nil)(), "must not be nil")
	})

	t.Run("ko because of empty", func(st *testing.T) {
		require.ErrorContains(st, MapRequired("agents", map[string]string{})(), "must not be empty")
	})
}
