package idx

import (
	"testing"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/require"
)

func TestBeforeTime(t *testing.T) {
	now := time.Now()

	id := BeforeTime(now)
	uid, err := ksuid.Parse(id)
	require.Nil(t, err)
	require.GreaterOrEqual(t, now.UnixMilli(), uid.Time().UnixMilli())
}

func TestAfterTime(t *testing.T) {
	now := time.Now()

	id := AfterTime(now)
	uid, err := ksuid.Parse(id)
	require.Nil(t, err)
	require.LessOrEqual(t, now.UnixMilli(), uid.Time().UnixMilli())
}
