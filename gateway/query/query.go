package query

import (
	"time"
)

var (
	SearchMaxChar = 256
	SearchMinChar = 3
	LimitMin      = 5
	LimitMax      = 100
	SizeMin       = 10
	SizeMax       = 500
	PageMin       = 1
	PageMax       = 100
	IdsMax        = 100

	DefaultPagingQuery = &Query{
		Limit: LimitMin,
		Page:  PageMin,
		Size:  SizeMin,
	}
)

// Query is an object that holds the query parameters for a listing
type Query struct {
	Search string
	// paging
	Limit int
	Page  int
	Ids   []string
	// scanning
	Cursor string
	Size   int
	From   time.Time
	To     time.Time
}

func (query *Query) Clone() *Query {
	return &Query{
		Search: query.Search,
		Limit:  query.Limit,
		Page:   query.Page,
		Ids:    query.Ids,
		Cursor: query.Cursor,
		Size:   query.Size,
		From:   query.From,
		To:     query.To,
	}
}
