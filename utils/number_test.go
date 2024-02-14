package utils

import (
	"math"
	"testing"

	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

var (
	min = testdata.Fake.IntBetween(0, math.MaxInt-1)
	max = min + 1
)

func TestMin(t *testing.T) {
	require.Equal(t, Min(min, max), min)
	require.Equal(t, Min(max, min), min)
}

func TestMax(t *testing.T) {
	require.Equal(t, Max(min, max), max)
	require.Equal(t, Max(max, min), max)
}
