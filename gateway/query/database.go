package query

import "github.com/kanthorlabs/common/persistence/database"

func (query *Query) ToDbPagingQuery() *database.PagingQuery {
	return &database.PagingQuery{
		Search: query.Search,
		Limit:  query.Limit,
		Page:   query.Page,
		Ids:    query.Ids,
	}
}
