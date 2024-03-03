package database

import (
	"strings"
	"testing"

	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/testify"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestPagingQuery_Clone(t *testing.T) {
	original := DefaultPagingQuery

	clone := original.Clone()
	clone.Search = testdata.Fake.Address().City()
	require.NotEqual(t, original.Search, clone.Search)
}

func TestPagingQuery_Sqlx(t *testing.T) {
	db := testify.GormStart(t)
	defer testify.GormEnd(t, db)

	records, ids := testify.GormInsert(t, db, testdata.Fake.IntBetween(10, 100))

	t.Run("with ids", func(st *testing.T) {
		query := DefaultPagingQuery.Clone()
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

	t.Run("with search keyword", func(st *testing.T) {
		query := DefaultPagingQuery.Clone()
		record := records[lo.Sample(ids)]
		query.Search = record.Username[:strings.Index(record.Username, "/")]

		var rows []*testdata.User
		tx := query.Sqlx(db.Model(&testdata.User{}), "id", []string{"username"}).Find(&rows)
		require.NoError(st, tx.Error)
		require.Equal(st, 1, len(rows))

		require.Equal(st, record, rows[0])
	})
}

func TestPagingQuery_SqlxCount(t *testing.T) {
	db := testify.GormStart(t)
	defer testify.GormEnd(t, db)

	records, ids := testify.GormInsert(t, db, testdata.Fake.IntBetween(10, 100))

	t.Run("with ids", func(st *testing.T) {
		query := DefaultPagingQuery.Clone()
		// searching is subset of ids
		searching := ids[:testdata.Fake.IntBetween(1, len(ids)-1)]
		query.Ids = searching

		var count int64
		tx := query.SqlxCount(db.Model(&testdata.User{}), "id", []string{"username"}).Count(&count)
		require.NoError(st, tx.Error)
		require.Equal(st, int64(len(query.Ids)), count)
	})

	t.Run("with search keyword", func(st *testing.T) {
		query := DefaultPagingQuery.Clone()
		record := records[lo.Sample(ids)]
		query.Search = record.Username[:strings.Index(record.Username, "/")]

		var count int64
		tx := query.SqlxCount(db.Model(&testdata.User{}), "id", []string{"username"}).Count(&count)
		require.NoError(st, tx.Error)
		require.Equal(st, int64(1), count)
	})
}
