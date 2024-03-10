package query

import (
	"net/http"
	"strconv"
	"time"

	"github.com/kanthorlabs/common/utils"
)

func FromHttpx(r *http.Request) *Query {
	query := &Query{Ids: []string{}}
	query.Search = r.URL.Query().Get("_q")

	// paging
	query.Limit = HttpxNumber(r, "_limit", LimitMin, LimitMax)
	query.Page = HttpxNumber(r, "_page", PageMin, PageMax)
	if len(r.URL.Query()["_ids"]) > 0 {
		query.Ids = r.URL.Query()["_ids"]
	}

	// scaning
	query.Cursor = r.URL.Query().Get("_cursor")
	query.Size = HttpxNumber(r, "_size", SizeMin, SizeMax)

	query.From = time.Now().UTC()
	from := HttpxNumber(r, "_from", int64(0), query.From.UnixMilli())
	if from > 0 {
		query.From = time.UnixMilli(from)
	}

	query.To = time.Now().UTC().Add(time.Minute)
	to := HttpxNumber(r, "_to", int64(0), query.To.UnixMilli())
	if to > 0 {
		query.To = time.UnixMilli(to)
	}
	return query
}

func HttpxNumber[T int | int64](r *http.Request, name string, min, max T) (value T) {
	if r.URL.Query().Get(name) != "" {
		// ignore error, we will use the default value in range of min and max
		v, _ := strconv.ParseInt(r.URL.Query().Get(name), 10, 0)
		value = T(v)
	}

	value = utils.Min(utils.Max(value, min), max)
	return
}
