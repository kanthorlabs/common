package query

import "github.com/kanthorlabs/common/persistence/datastore"

func (query *Query) ToDsScanningQuery() *datastore.ScanningQuery {
	return &datastore.ScanningQuery{
		Search: query.Search,
		Cursor: query.Cursor,
		Size:   query.Size,
		From:   query.From,
		To:     query.To,
	}
}
