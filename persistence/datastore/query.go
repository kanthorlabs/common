package datastore

import (
	"fmt"
	"time"

	"github.com/kanthorlabs/common/idx"
	"gorm.io/gorm"
)

var (
	SearchMaxChar = 256
	SearchMinChar = 3
	SizeMin       = 5
	SizeMax       = 100
)

type ScanningCondition struct {
	PrimaryKeyNs  string
	PrimaryKeyCol string
}

type ScanningQuery struct {
	Cursor string
	Search string
	Size   int
	From   time.Time
	To     time.Time
}

func (query *ScanningQuery) Clone() *ScanningQuery {
	return &ScanningQuery{
		Cursor: query.Cursor,
		Search: query.Search,
		Size:   query.Size,
		From:   query.From,
		To:     query.To,
	}
}

func (query *ScanningQuery) Sqlx(tx *gorm.DB, condition *ScanningCondition) *gorm.DB {
	low := idx.Build(condition.PrimaryKeyNs, idx.BeforeTime(query.From))
	high := idx.Build(condition.PrimaryKeyNs, idx.AfterTime(query.To))

	tx = tx.
		Where(fmt.Sprintf(`%s > ?`, condition.PrimaryKeyCol), low).
		Where(fmt.Sprintf(`%s < ?`, condition.PrimaryKeyCol), high).
		Order(fmt.Sprintf(`%s DESC`, condition.PrimaryKeyCol)).
		Limit(query.Size)

	if query.Search != "" {
		// IMPORTANT: only support search by primary key
		// our primary key is often conbined from multiple columns
		// so you can search with the second column of the primary key
		// when and only when you added the first column to the where condition
		// for example, your primary key is message_pk(app_id, id)
		// you can only match the where condition for the ID column
		// when you add the where condition for the app_id column before
		// message: app_id = ? AND id = ?
		tx = tx.Where(fmt.Sprintf(`%s = ?`, condition.PrimaryKeyCol), query.Search)
	}

	if query.Cursor != "" {
		tx = tx.Where(fmt.Sprintf(`%s < ?`, condition.PrimaryKeyCol), query.Cursor)
	}

	return tx
}

type ScanningRecord[T any] struct {
	Data  T
	Error error
}
