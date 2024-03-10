package database

import (
	"fmt"

	"github.com/kanthorlabs/common/utils"
	"github.com/kanthorlabs/common/validator"
	"gorm.io/gorm"
)

var (
	SearchMaxChar = 256
	SearchMinChar = 3
	LimitMin      = 5
	LimitMax      = 100
	PageMin       = 1
	PageMax       = 100

	DefaultPagingQuery = &PagingQuery{Limit: LimitMin, Page: PageMin}
)

type PagingQuery struct {
	Search string
	Limit  int
	Page   int
	Ids    []string
}

func (query *PagingQuery) Clone() *PagingQuery {
	return &PagingQuery{
		Search: query.Search,
		Limit:  query.Limit,
		Page:   query.Page,
		Ids:    query.Ids,
	}
}

func (query *PagingQuery) Validate() error {
	return validator.Validate(
		validator.StringLenIfNotEmpty("DATABASE.QUERY.SEARCH", query.Search, SearchMinChar, SearchMaxChar),
		validator.NumberInRange("DATABASE.QUERY.LIMIT", query.Limit, LimitMin, LimitMax),
		validator.NumberInRange("DATABASE.QUERY.PAGE", query.Page, PageMin, PageMax),
		// only allow retrieving a maximum of 100 records at a time as limit
		validator.SliceMaxLength("DATABASE.QUERY.IDS", query.Ids, LimitMax),
		validator.Slice(query.Ids, func(i int, item *string) error {
			key := fmt.Sprintf("DATABASE.QUERY.IDS[%d]", i)
			return validator.StringLen(key, *item, SearchMinChar, SearchMaxChar)()
		}),
	)
}

func (query *PagingQuery) Sqlx(tx *gorm.DB, primaryCol string, searchCols []string) *gorm.DB {
	tx = query.sqlx(tx, primaryCol, searchCols)

	if len(query.Ids) == 0 {
		offset := utils.Max((query.Page-1)*query.Limit, 0)
		tx = tx.Limit(query.Limit).Offset(offset)
	}

	return tx
}

func (query *PagingQuery) SqlxCount(tx *gorm.DB, primaryCol string, searchCols []string) *gorm.DB {
	return query.sqlx(tx, primaryCol, searchCols)
}

func (query *PagingQuery) sqlx(tx *gorm.DB, primaryCol string, searchCols []string) *gorm.DB {
	if len(query.Ids) > 0 {
		return tx.Where(fmt.Sprintf("%s IN ?", primaryCol), query.Ids)
	}

	for i := range searchCols {
		tx = tx.Where(fmt.Sprintf(`%s LIKE ?`, searchCols[i]), query.Search+"%")
	}

	return tx
}
