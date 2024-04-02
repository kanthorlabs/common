package webhook

import (
	"testing"
	"time"

	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestOptions(t *testing.T) {
	options := &VerifyOptions{}
	require.False(t, options.TimestampToleranceIgnore)
	require.Equal(t, time.Duration(0), options.TimestampToleranceDuration)

	duration := time.Since(testdata.Fake.Time().TimeBetween(time.Now().Add(-time.Hour*2), time.Now().Add(-time.Hour)))
	TimestampToleranceDuration(duration)(options)
	require.Equal(t, duration, options.TimestampToleranceDuration)

	TimestampToleranceIgnore()(options)
	require.True(t, options.TimestampToleranceIgnore)

	defer func() {
		if r := recover(); r != nil {
			require.ErrorIs(t, r.(error), ErrTimestampToleranceDurationTooSmall)
		}
	}()
}
