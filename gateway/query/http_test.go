package query

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func Test_FromHttpx(t *testing.T) {
	t.Run("OK ", func(st *testing.T) {
		q := url.Values{}
		q.Add("_q", testquery.Search)
		q.Add("_limit", strconv.Itoa(testquery.Limit))
		q.Add("_page", strconv.Itoa(testquery.Page))
		for _, id := range testquery.Ids {
			q.Add("_ids", id)
		}
		q.Add("_cursor", testquery.Cursor)
		q.Add("_size", strconv.Itoa(testquery.Size))
		q.Add("_from", fmt.Sprintf("%d", testquery.From.UnixMilli()))
		q.Add("_to", fmt.Sprintf("%d", testquery.To.UnixMilli()))

		u := url.URL{}
		u.RawQuery = q.Encode()

		req, err := http.NewRequest(http.MethodGet, u.String(), nil)
		require.NoError(st, err)

		query := FromHttpx(req)
		require.Equal(st, testquery.Search, query.Search)
		require.Equal(st, testquery.Limit, query.Limit)
		require.Equal(st, testquery.Page, query.Page)
		require.Equal(st, testquery.Ids, query.Ids)
		require.Equal(st, testquery.Cursor, query.Cursor)
		require.Equal(st, testquery.Size, query.Size)
		// datetime is not equal because of the precision in nanoseconds
		require.Equal(st, testquery.From.UnixMilli(), query.From.UnixMilli())
		require.Equal(st, testquery.To.UnixMilli(), query.To.UnixMilli())
	})

	t.Run("KO - default value", func(st *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "", nil)
		require.NoError(st, err)

		query := FromHttpx(req)
		require.Equal(st, "", query.Search)
		require.Equal(st, LimitMin, query.Limit)
		require.Equal(st, PageMin, query.Page)
		require.Equal(st, []string{}, query.Ids)
		require.Equal(st, "", query.Cursor)
		require.Equal(st, SizeMin, query.Size)
		require.GreaterOrEqual(st, query.From.UnixMilli(), int64(0))
		require.GreaterOrEqual(st, query.To.UnixMilli(), int64(0))
	})

	t.Run("OK - parse _limit and _page error", func(st *testing.T) {
		q := url.Values{}
		q.Add("_limit", testdata.Fake.RandomLetter())
		q.Add("_page", testdata.Fake.RandomLetter())

		u := url.URL{}
		u.RawQuery = q.Encode()

		req, err := http.NewRequest(http.MethodGet, u.String(), nil)
		require.NoError(st, err)

		query := FromHttpx(req)
		require.Equal(st, "", query.Search)
		require.Equal(st, LimitMin, query.Limit)
		require.Equal(st, PageMin, query.Page)
		require.Equal(st, []string{}, query.Ids)
		require.Equal(st, "", query.Cursor)
		require.Equal(st, SizeMin, query.Size)
		require.GreaterOrEqual(st, query.From.UnixMilli(), int64(0))
		require.GreaterOrEqual(st, query.To.UnixMilli(), int64(0))
	})

	t.Run("OK - parse _size error", func(st *testing.T) {
		q := url.Values{}
		q.Add("_size", testdata.Fake.RandomLetter())

		u := url.URL{}
		u.RawQuery = q.Encode()

		req, err := http.NewRequest(http.MethodGet, u.String(), nil)
		require.NoError(st, err)

		query := FromHttpx(req)
		require.Equal(st, "", query.Search)
		require.Equal(st, LimitMin, query.Limit)
		require.Equal(st, PageMin, query.Page)
		require.Equal(st, []string{}, query.Ids)
		require.Equal(st, "", query.Cursor)
		require.Equal(st, SizeMin, query.Size)
		require.GreaterOrEqual(st, query.From.UnixMilli(), int64(0))
		require.GreaterOrEqual(st, query.To.UnixMilli(), int64(0))
	})

	t.Run("OK - parse _from and _to error", func(st *testing.T) {
		q := url.Values{}
		q.Add("_from", testdata.Fake.RandomLetter())
		q.Add("_to", testdata.Fake.RandomLetter())

		u := url.URL{}
		u.RawQuery = q.Encode()

		req, err := http.NewRequest(http.MethodGet, u.String(), nil)
		require.NoError(st, err)

		query := FromHttpx(req)
		require.Equal(st, "", query.Search)
		require.Equal(st, LimitMin, query.Limit)
		require.Equal(st, PageMin, query.Page)
		require.Equal(st, []string{}, query.Ids)
		require.Equal(st, "", query.Cursor)
		require.Equal(st, SizeMin, query.Size)
		require.GreaterOrEqual(st, query.From.UnixMilli(), int64(0))
		require.GreaterOrEqual(st, query.To.UnixMilli(), int64(0))
	})
}

func Test_HttpxNumber(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "", nil)
		require.NoError(st, err)

		require.Equal(st, int64(0), HttpxNumber(req, "test", int64(0), int64(10)))
	})

	t.Run("OK - parse error", func(st *testing.T) {
		q := url.Values{}
		q.Add("test", testdata.Fake.RandomLetter())

		u := url.URL{}
		u.RawQuery = q.Encode()

		req, err := http.NewRequest(http.MethodGet, u.String(), nil)
		require.NoError(st, err)

		require.Equal(st, int64(0), HttpxNumber(req, "test", int64(0), int64(10)))
	})
}
