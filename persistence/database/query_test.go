package database

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
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

func TestPagingQuery_Validate(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		query := DefaultPagingQuery.Clone()

		require.NoError(st, query.Validate())
	})

	t.Run("KO - search error", func(st *testing.T) {
		query := DefaultPagingQuery.Clone()

		query.Search = testdata.Fake.Lorem().Sentence(SearchMaxChar + 1)
		require.ErrorContains(st, query.Validate(), "DATABASE.QUERY.SEARCH")
	})

	t.Run("KO - limit error", func(st *testing.T) {
		query := DefaultPagingQuery.Clone()

		query.Limit = LimitMax + 1
		require.ErrorContains(st, query.Validate(), "DATABASE.QUERY.LIMIT")
	})

	t.Run("KO - page error", func(st *testing.T) {
		query := DefaultPagingQuery.Clone()

		query.Page = PageMax + 1
		require.ErrorContains(st, query.Validate(), "DATABASE.QUERY.PAGE")
	})

	t.Run("KO - ids error", func(st *testing.T) {
		query := DefaultPagingQuery.Clone()

		query.Ids = append(query.Ids, testdata.Fake.Lorem().Sentence(SearchMaxChar+1))
		require.ErrorContains(st, query.Validate(), fmt.Sprintf("DATABASE.QUERY.IDS[%d]", len(query.Ids)-1))

		for i := 0; i < LimitMax+1; i++ {
			query.Ids = append(query.Ids, uuid.NewString())
		}
		require.ErrorContains(st, query.Validate(), "DATABASE.QUERY.IDS")
	})
}

func TestPagingQuery_Sqlx(t *testing.T) {
	db := testify.GormStart(t)
	defer testify.GormEnd(t, db)

	records, ids := testify.GormInsert(t, db, testdata.Fake.IntBetween(10, 100))

	t.Run("OK", func(st *testing.T) {
		query := DefaultPagingQuery.Clone()

		var rows []*testdata.User
		tx := query.Sqlx(db.Model(&testdata.User{}), "id", []string{"username"}).Find(&rows)
		require.NoError(st, tx.Error)
		require.Equal(st, query.Limit, len(rows))
	})

	t.Run("OK - with ids", func(st *testing.T) {
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

	t.Run("OK - with search keyword", func(st *testing.T) {
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

	t.Run("OK", func(st *testing.T) {
		query := DefaultPagingQuery.Clone()

		var count int64
		tx := query.SqlxCount(db.Model(&testdata.User{}), "id", []string{"username"}).Count(&count)
		require.NoError(st, tx.Error)
		require.Equal(st, int64(len(records)), count)
	})

	t.Run("OK - with ids", func(st *testing.T) {
		query := DefaultPagingQuery.Clone()
		// searching is subset of ids
		searching := ids[:testdata.Fake.IntBetween(1, len(ids)-1)]
		query.Ids = searching

		var count int64
		tx := query.SqlxCount(db.Model(&testdata.User{}), "id", []string{"username"}).Count(&count)
		require.NoError(st, tx.Error)
		require.Equal(st, int64(len(query.Ids)), count)
	})

	t.Run("OK - with search keyword", func(st *testing.T) {
		query := DefaultPagingQuery.Clone()
		record := records[lo.Sample(ids)]
		query.Search = record.Username[:strings.Index(record.Username, "/")]

		var count int64
		tx := query.SqlxCount(db.Model(&testdata.User{}), "id", []string{"username"}).Count(&count)
		require.NoError(st, tx.Error)
		require.Equal(st, int64(1), count)
	})
}
