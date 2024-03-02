package validator

import (
	"testing"

	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestSliceRequired(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		agents := []string{
			testdata.Fake.UserAgent().Chrome(),
			testdata.Fake.UserAgent().Firefox(),
			testdata.Fake.UserAgent().Safari(),
		}
		require.NoError(st, SliceRequired("agents", agents)())
	})

	t.Run("KO - nil slice", func(st *testing.T) {
		require.ErrorContains(st, SliceRequired[string]("agents", nil)(), "must not be nil")
	})

	t.Run("KO - empty list", func(st *testing.T) {
		require.ErrorContains(st, SliceRequired[string]("agents", []string{})(), "must not be empty")
	})
}

func TestSliceMaxLength(t *testing.T) {
	count := testdata.Fake.IntBetween(100000, 999999)
	items := make([]int, count)

	t.Run("OK", func(st *testing.T) {
		require.NoError(st, SliceMaxLength("items", items, count)())
	})

	t.Run("KO - exceeded maximum capacity", func(st *testing.T) {
		require.ErrorContains(st, SliceMaxLength("items", items, count-1)(), "is exceeded maximum capacity")
	})
}

func TestMapRequired(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		agents := map[string]string{
			"chrome":  testdata.Fake.UserAgent().Chrome(),
			"firefox": testdata.Fake.UserAgent().Firefox(),
			"safar":   testdata.Fake.UserAgent().Safari(),
		}
		require.NoError(st, MapRequired("agents", agents)())
	})

	t.Run("KO - nil map", func(st *testing.T) {
		require.ErrorContains(st, MapRequired[string, string]("agents", nil)(), "must not be nil")
	})

	t.Run("KO - empty map", func(st *testing.T) {
		require.ErrorContains(st, MapRequired("agents", map[string]string{})(), "must not be empty")
	})
}
