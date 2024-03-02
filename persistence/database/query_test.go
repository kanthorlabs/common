package database

import (
	"strings"
	"testing"

	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/testify"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestPagingQuery(t *testing.T) {
	original := DefaultPagingQuery

	t.Run(".Clone", func(st *testing.T) {
		clone := original.Clone()
		clone.Search = testdata.Fake.Address().City()

		require.NotEqual(st, original.Search, clone.Search)
	})

	t.Run(".Sqlx/Ids", func(st *testing.T) {
		db := testify.GormStart(st)
		defer testify.GormEnd(st, db)

		records, ids := testify.GormInsert(st, db, testdata.Fake.IntBetween(10, 100))

		query := original.Clone()
		// searching is subset of ids
		searching := ids[:testdata.Fake.IntBetween(1, len(ids)-1)]
		query.Ids = searching

		var rows []*testdata.User
		tx := query.Sqlx(db.Model(&testdata.User{}), "id", []string{"username"}).Find(&rows)
		require.NoError(st, tx.Error)
		require.Equal(st, len(query.Ids), len(rows))

		for i := range rows {
			require.Equal(st, rows[i], records[rows[i].Id])
		}
	})

	t.Run(".Sqlx/Search", func(st *testing.T) {
		db := testify.GormStart(st)
		defer testify.GormEnd(st, db)

		records, ids := testify.GormInsert(st, db, testdata.Fake.IntBetween(10, 100))

		query := original.Clone()
		record := records[lo.Sample(ids)]
		query.Search = record.Username[:strings.Index(record.Username, "/")]

		var rows []*testdata.User
		tx := query.Sqlx(db.Model(&testdata.User{}), "id", []string{"username"}).Find(&rows)
		require.NoError(st, tx.Error)
		require.Equal(st, 1, len(rows))

		require.Equal(st, record, rows[0])
	})

	t.Run(".SqlxCount/Ids", func(st *testing.T) {
		db := testify.GormStart(st)
		defer testify.GormEnd(st, db)

		_, ids := testify.GormInsert(st, db, testdata.Fake.IntBetween(10, 100))

		query := original.Clone()
		// searching is subset of ids
		searching := ids[:testdata.Fake.IntBetween(1, len(ids)-1)]
		query.Ids = searching

		var count int64
		tx := query.SqlxCount(db.Model(&testdata.User{}), "id", []string{"username"}).Count(&count)
		require.NoError(st, tx.Error)
		require.Equal(st, int64(len(query.Ids)), count)

	})

	t.Run(".SqlxCount/Search", func(st *testing.T) {
		db := testify.GormStart(st)
		defer testify.GormEnd(st, db)

		records, ids := testify.GormInsert(st, db, testdata.Fake.IntBetween(10, 100))

		query := original.Clone()
		record := records[lo.Sample(ids)]
		query.Search = record.Username[:strings.Index(record.Username, "/")]

		var count int64
		tx := query.SqlxCount(db.Model(&testdata.User{}), "id", []string{"username"}).Count(&count)
		require.NoError(st, tx.Error)
		require.Equal(st, int64(1), count)
	})

}
