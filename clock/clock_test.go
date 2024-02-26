package clock

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClock(t *testing.T) {
	clock := New()

	now := clock.Now()
	require.Equal(t, now.UnixMilli(), clock.UnixMilli(now.UnixMilli()).UnixMilli())
}
