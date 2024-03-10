package query

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

var testquery = &Query{
	Search: testdata.Fake.App().Name(),
	Limit:  testdata.Fake.IntBetween(LimitMin, LimitMax),
	Page:   testdata.Fake.IntBetween(PageMin, PageMax),
	Ids:    []string{uuid.New().String(), uuid.New().String()},
	Cursor: uuid.NewString(),
	Size:   testdata.Fake.IntBetween(SizeMin, SizeMax),
	From:   testdata.Fake.Time().TimeBetween(time.Now().Add(-time.Hour*24*14), time.Now().Add(-time.Hour*24*7)),
	To:     testdata.Fake.Time().TimeBetween(time.Now().Add(-time.Hour*24*7), time.Now()),
}

func TestQuery_Clone(t *testing.T) {
	clone := testquery.Clone()

	require.Equal(t, testquery.Search, clone.Search)
	require.Equal(t, testquery.Limit, clone.Limit)
	require.Equal(t, testquery.Page, clone.Page)
	require.Equal(t, testquery.Ids, clone.Ids)
	require.Equal(t, testquery.Cursor, clone.Cursor)
	require.Equal(t, testquery.Size, clone.Size)
	require.Equal(t, testquery.From, clone.From)
	require.Equal(t, testquery.To, clone.To)
}
