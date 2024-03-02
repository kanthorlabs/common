package datastore

import (
	"testing"

	"github.com/kanthorlabs/common/idx"
	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/testify"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestScanningQuery(t *testing.T) {
	condition := &ScanningCondition{
		PrimaryKeyNs:  testdata.UserNs,
		PrimaryKeyCol: "id",
	}
	original := &ScanningQuery{}

	t.Run(".Clone", func(st *testing.T) {
		clone := original.Clone()
		clone.Search = testdata.Fake.Address().City()

		require.NotEqual(st, original.Search, clone.Search)
	})

	t.Run(".Sqlx", func(st *testing.T) {
		db := testify.GormStart(st)
		defer testify.GormEnd(st, db)

		query := original.Clone()
		records, _, _, size, count := setup(st, db, query)

		// request more records than the range of .From and .To
		query.Size = count

		var rows []*testdata.User
		tx := query.Sqlx(db, condition).Find(&rows)

		require.NoError(st, tx.Error)
		// make sure we only retrieve records in the range of .From and .To
		require.Equal(st, size, len(rows))

		for i := range rows {
			require.Equal(st, rows[i], records[rows[i].Id])
		}
	})

	t.Run(".Sqlx/Search", func(st *testing.T) {
		db := testify.GormStart(st)
		defer testify.GormEnd(st, db)

		query := original.Clone()
		records, ids, mid, size, count := setup(st, db, query)
		// use id to search
		search := ids[count-mid-(size/2)]
		query.Search = search

		var rows []*testdata.User
		tx := query.Sqlx(db, condition).Find(&rows)

		require.NoError(st, tx.Error)
		require.Equal(st, 1, len(rows))

		require.Equal(st, rows[0], records[search])
	})

	t.Run(".Sqlx/Cursor", func(st *testing.T) {
		db := testify.GormStart(st)
		defer testify.GormEnd(st, db)

		query := original.Clone()
		records, ids, mid, size, count := setup(st, db, query)
		// use id to search
		cursor := ids[count-mid-(size/2)-1]
		query.Cursor = cursor

		var rows []*testdata.User
		tx := query.Sqlx(db, condition).Find(&rows)

		require.NoError(st, tx.Error)
		require.Equal(st, size/2, len(rows))

		for i := range rows {
			require.Equal(st, rows[i], records[rows[i].Id])
		}

		// make sure all rows is less than the cursor because ScanningOrderDesc is used by default
		require.Greater(st, rows[0].Id, rows[len(rows)-1].Id)
		require.Greater(st, cursor, rows[0].Id)
	})
}

func setup(t *testing.T, db *gorm.DB, query *ScanningQuery) (map[string]*testdata.User, []string, int, int, int) {
	mid := testdata.Fake.IntBetween(0, 9)
	size := testdata.Fake.IntBetween(11, 99)
	count := testdata.Fake.IntBetween(101, 101+size)
	records, ids := testify.GormInsert(t, db, count)

	query.Size = size

	from, err := idx.ToTime(ids[count-mid-1])
	require.NoError(t, err)
	query.From = from
	to, err := idx.ToTime(ids[count-mid-size])
	require.NoError(t, err)
	query.To = to

	return records, ids, mid, size, count
}
