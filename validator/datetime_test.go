package validator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDatetimeBefore(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		err := DatetimeBefore("prop", time.Now(), time.Now().Add(time.Hour))()
		require.NoError(st, err)
	})

	t.Run("KO", func(st *testing.T) {
		err := DatetimeBefore("prop", time.Now().Add(time.Hour), time.Now())()
		require.ErrorContains(st, err, "must before")
	})

	t.Run("KO - value is before MinDatetime", func(st *testing.T) {
		err := DatetimeBefore("prop", MinDatetime.Add(-time.Hour), time.Now())()
		require.ErrorContains(st, err, "must after "+MinDatetime.Format(time.RFC3339Nano))
	})

	t.Run("KO - target is before MinDatetime", func(st *testing.T) {
		err := DatetimeBefore("prop", time.Now(), MinDatetime.Add(-time.Hour))()
		require.ErrorContains(st, err, "must after "+MinDatetime.Format(time.RFC3339Nano))
	})

}
