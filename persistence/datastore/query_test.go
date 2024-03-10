package datastore

import (
	"testing"
	"time"

	"github.com/kanthorlabs/common/idx"
	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/testify"
	"github.com/kanthorlabs/common/validator"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestScanningQuery_Clone(t *testing.T) {
	original := &ScanningQuery{}

	clone := original.Clone()
	clone.Search = testdata.Fake.Address().City()
	require.NotEqual(t, original.Search, clone.Search)
}

func TestScanningQuery_Validate(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		query := &ScanningQuery{
			Size: SizeMin,
			From: validator.MinDatetime,
			To:   time.Now(),
		}

		require.NoError(st, query.Validate())
	})

	t.Run("KO - search error", func(st *testing.T) {
		query := &ScanningQuery{
			Size: SizeMin,
			From: validator.MinDatetime,
			To:   time.Now(),
		}
		query.Search = testdata.Fake.Lorem().Sentence(SearchMaxChar + 1)

		require.ErrorContains(st, query.Validate(), "DATASTORE.QUERY.SEARCH")
	})

	t.Run("KO - size error", func(st *testing.T) {
		query := &ScanningQuery{
			Size: SizeMin,
			From: validator.MinDatetime,
			To:   time.Now(),
		}
		query.Size = SizeMax + 1

		require.ErrorContains(st, query.Validate(), "DATASTORE.QUERY.SIZE")
	})

	t.Run("KO - from error", func(st *testing.T) {
		query := &ScanningQuery{
			Size: SizeMin,
			To:   time.Now(),
		}

		require.ErrorContains(st, query.Validate(), "DATASTORE.QUERY.FROM")
	})

	t.Run("KO - to error", func(st *testing.T) {
		query := &ScanningQuery{
			Size: SizeMin,
			From: validator.MinDatetime,
		}

		require.ErrorContains(st, query.Validate(), "DATASTORE.QUERY.FROM")
	})
}

func TestScanningQuery_Sqlx(t *testing.T) {
	condition := &ScanningCondition{
		PrimaryKeyNs:  testdata.UserNs,
		PrimaryKeyCol: "id",
	}
	original := &ScanningQuery{}

	db := testify.GormStart(t)
	defer testify.GormEnd(t, db)

	t.Run("OK", func(st *testing.T) {
		query := original.Clone()
		records, _, _, size, count := setup(t, db, query)

		// request more records than the range of .From and .To
		query.Size = count

		var rows []*testdata.User
		tx := query.Sqlx(db, condition).Find(&rows)

		require.NoError(t, tx.Error)
		// make sure we only retrieve records in the range of .From and .To
		require.Equal(t, size, len(rows))

		for i := range rows {
			require.Equal(t, rows[i], records[rows[i].Id])
		}
	})

	t.Run("OK - with search keyword", func(st *testing.T) {
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

	t.Run("OK - with cursor", func(st *testing.T) {
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
	// .Where("1=1") is a hack to bypass the gorm's check for empty condition
	require.NoError(t, db.Model(&testdata.User{}).Where("1=1").Delete(nil).Error)

	mid := testdata.Fake.IntBetween(0, 9)
	size := testdata.Fake.IntBetween(99, 999)
	count := testdata.Fake.IntBetween(1001, 1001+size)
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
