package idx

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/testdata"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	ns := testdata.Fake.Color().SafeColorName()
	require.True(t, strings.HasPrefix(New(ns), ns))
}

func TestBuild(t *testing.T) {
	ns := testdata.Fake.Color().SafeColorName()
	id := ksuid.New().String()
	require.Equal(t, fmt.Sprintf("%s_%s", ns, id), Build(ns, id))
}

func TestFromTime(t *testing.T) {
	ns := testdata.Fake.Color().SafeColorName()
	require.True(t, strings.HasPrefix(FromTime(ns, time.Now()), ns))
}

func TestToTime(t *testing.T) {
	t.Run("ok", func(st *testing.T) {
		now := time.Now()
		id := FromTime(testdata.Fake.Color().SafeColorName(), now)

		ts, err := ToTime(id)
		require.Nil(st, err)

		// idx is only guaranteed to be unique at the second level
		require.Equal(st, now.Unix(), ts.Unix())
	})

	t.Run("ko because of malformed format", func(st *testing.T) {
		_, err := ToTime(uuid.NewString())
		require.ErrorContains(st, err, "IDX.MALFORMED_FORMAT.ERROR")
	})

	t.Run("ko because of parsing error", func(st *testing.T) {
		_, err := ToTime(fmt.Sprintf("u_%s", uuid.NewString()))
		require.ErrorContains(st, err, "IDX.PARSE.ERROR")
	})
}
