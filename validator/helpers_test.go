package validator

import (
	"testing"

	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

var (
	low  = testdata.Fake.IntBetween(1, 100)
	high = testdata.Fake.IntBetween(100000, 1000000)
	step = testdata.Fake.IntBetween(100, 1000)
)

func svalidate(x int) func(i int, item *int) error {
	return func(i int, item *int) error {
		if *item >= x {
			return testdata.ErrorGeneric
		}
		return nil
	}
}
func TestSlice(t *testing.T) {
	items := []int{}
	for i := low; i < high; i += step {
		items = append(items, i)
	}
	t.Run("OK", func(st *testing.T) {
		err := Validate(Slice(items, svalidate(high+1)))
		require.NoError(st, err)
	})

	t.Run("KO - item error", func(st *testing.T) {
		err := Validate(Slice(items, svalidate(high-step)))
		require.ErrorIs(st, err, testdata.ErrorGeneric)
	})
}

func mvalidate(x int) func(refId string, item int) error {
	return func(refId string, item int) error {
		if item >= x {
			return testdata.ErrorGeneric
		}
		return nil
	}
}

func TestMap(t *testing.T) {
	items := map[string]int{}
	for i := low; i < high; i += step {
		items[testdata.Fake.UUID().V4()] = i
	}
	t.Run("OK", func(st *testing.T) {
		err := Validate(Map(items, mvalidate(high+1)))
		require.NoError(st, err)
	})

	t.Run("KO - item error", func(st *testing.T) {
		err := Validate(Map(items, mvalidate(high-step)))
		require.ErrorIs(st, err, testdata.ErrorGeneric)
	})
}
