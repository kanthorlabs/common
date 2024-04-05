package webhook

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/kanthorlabs/common/cipher/signature"
	"github.com/kanthorlabs/common/idx"
	"github.com/kanthorlabs/common/testdata"
	"github.com/kanthorlabs/common/utils"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		wh, err := New(keys)
		require.NoError(st, err)
		require.NotNil(st, wh)
	})

	t.Run("OK - custom key namespace", func(st *testing.T) {
		ns := testdata.Fake.RandomStringWithLength(5)
		keys := []string{
			idx.Build(ns, utils.RandomString(128)),
		}

		wh, err := New(keys, KeyNamespace(ns))
		require.NoError(st, err)
		require.NotNil(st, wh)
	})

	t.Run("KO - no keys error", func(st *testing.T) {
		_, err := New([]string{})
		require.ErrorContains(st, err, "must not be empty")
	})

	t.Run("KO - too many keys error", func(st *testing.T) {
		keys := make([]string, MaxKeys+1)
		for i := range keys {
			keys[i] = idx.Build(DefaultKeyNs, utils.RandomString(128))
		}
		_, err := New(keys)
		require.ErrorContains(st, err, "is exceeded maximum capacity")
	})

	t.Run("KO - key format error", func(st *testing.T) {
		_, err := New([]string{utils.RandomString(128)})
		require.ErrorContains(st, err, " must be started with")
	})
}

func TestSign(t *testing.T) {
	wh, err := New(keys)
	require.NoError(t, err)

	t.Run("OK", func(st *testing.T) {
		id := idx.New("msg")
		ts := fmt.Sprintf("%d", time.Now().UnixMilli())
		body := utils.Stringify(map[string]any{
			"app_id": "msg_2ePVr2tTfiJA20mN8wkc8EkGZu4",
			"type":   "testing.openapi",
			"object": map[string]any{"from_client": "openapi", "say": "hello"},
		})

		signatures := wh.Sign(id, ts, body)
		require.Equal(st, len(keys), len(signatures))
	})
}

func TestVerify(t *testing.T) {
	wh, err := New(keys)
	require.NoError(t, err)

	body := utils.Stringify(map[string]any{
		"app_id": "msg_2ePVr2tTfiJA20mN8wkc8EkGZu4",
		"type":   "testing.openapi",
		"object": map[string]any{"from_client": "openapi", "say": "hello"},
	})

	t.Run("OK", func(st *testing.T) {
		id := idx.New("msg")
		ts := fmt.Sprintf("%d", time.Now().UnixMilli())

		req := httptest.NewRequest(http.MethodPost, "/webhook/demo", io.NopCloser(strings.NewReader(body)))
		req.Header.Set(HeaderId, id)
		req.Header.Set(HeaderTimestamp, ts)

		key := testdata.Fake.RandomStringElement(keys)
		signature := signature.Sign(key, fmt.Sprintf("%s.%s.%s", id, ts, body))
		req.Header.Set(HeaderSignature, signature)

		require.NoError(st, wh.Verify(req, TimestampToleranceDuration(time.Hour)))
	})

	t.Run("OK - ignore timestamp check", func(st *testing.T) {
		id := idx.New("msg")
		ts := fmt.Sprintf("%d", time.Now().Add(-ToleranceDurationDefault*2).UnixMilli())

		req := httptest.NewRequest(http.MethodPost, "/webhook/demo", io.NopCloser(strings.NewReader(body)))
		req.Header.Set(HeaderId, id)
		req.Header.Set(HeaderTimestamp, ts)

		key := testdata.Fake.RandomStringElement(keys)
		signature := signature.Sign(key, fmt.Sprintf("%s.%s.%s", id, ts, body))
		req.Header.Set(HeaderSignature, signature)

		require.NoError(st, wh.Verify(req, TimestampToleranceIgnore()))
	})

	t.Run("OK - with longer timestamp check", func(st *testing.T) {
		id := idx.New("msg")
		ts := fmt.Sprintf("%d", time.Now().Add(-ToleranceDurationDefault*3).UnixMilli())

		req := httptest.NewRequest(http.MethodPost, "/webhook/demo", io.NopCloser(strings.NewReader(body)))
		req.Header.Set(HeaderId, id)
		req.Header.Set(HeaderTimestamp, ts)

		key := testdata.Fake.RandomStringElement(keys)
		signature := signature.Sign(key, fmt.Sprintf("%s.%s.%s", id, ts, body))
		req.Header.Set(HeaderSignature, signature)

		require.NoError(st, wh.Verify(req, TimestampToleranceDuration(ToleranceDurationDefault*3)))
	})

	t.Run("KO - read body error", func(st *testing.T) {
		id := idx.New("msg")
		ts := fmt.Sprintf("%d", time.Now().UnixMilli())

		req := httptest.NewRequest(http.MethodPost, "/webhook/demo", io.NopCloser(&testdata.BrokenReader{}))
		req.Header.Set(HeaderId, id)
		req.Header.Set(HeaderTimestamp, ts)

		key := testdata.Fake.RandomStringElement(keys)
		signature := signature.Sign(key, fmt.Sprintf("%s.%s.%s", id, ts, body))
		req.Header.Set(HeaderSignature, signature)

		require.ErrorIs(st, wh.Verify(req), testdata.ErrGeneric)
	})

	t.Run("KO - parse timestamp error", func(st *testing.T) {
		id := idx.New("msg")
		ts := fmt.Sprintf("xxx-%d", time.Now().UnixMilli())

		req := httptest.NewRequest(http.MethodPost, "/webhook/demo", io.NopCloser(strings.NewReader(body)))
		req.Header.Set(HeaderId, id)
		req.Header.Set(HeaderTimestamp, ts)

		key := testdata.Fake.RandomStringElement(keys)
		signature := signature.Sign(key, fmt.Sprintf("%s.%s.%s", id, ts, body))
		req.Header.Set(HeaderSignature, signature)

		require.ErrorContains(st, wh.Verify(req), "WEBHOOK.MESSAGE.TIMESTAMP")
	})

	t.Run("KO - timestamp is too old error", func(st *testing.T) {
		id := idx.New("msg")
		ts := fmt.Sprintf("%d", time.Now().Add(-ToleranceDurationDefault-time.Hour).UnixMilli())

		req := httptest.NewRequest(http.MethodPost, "/webhook/demo", io.NopCloser(strings.NewReader(body)))
		req.Header.Set(HeaderId, id)
		req.Header.Set(HeaderTimestamp, ts)

		key := testdata.Fake.RandomStringElement(keys)
		signature := signature.Sign(key, fmt.Sprintf("%s.%s.%s", id, ts, body))
		req.Header.Set(HeaderSignature, signature)

		require.ErrorIs(st, wh.Verify(req), ErrMessageTimestampTooOld)
	})

	t.Run("KO - timestamp is too new error", func(st *testing.T) {
		id := idx.New("msg")
		ts := fmt.Sprintf("%d", time.Now().Add(ToleranceDurationDefault+time.Hour).UnixMilli())

		req := httptest.NewRequest(http.MethodPost, "/webhook/demo", io.NopCloser(strings.NewReader(body)))
		req.Header.Set(HeaderId, id)
		req.Header.Set(HeaderTimestamp, ts)

		key := testdata.Fake.RandomStringElement(keys)
		signature := signature.Sign(key, fmt.Sprintf("%s.%s.%s", id, ts, body))
		req.Header.Set(HeaderSignature, signature)

		require.ErrorIs(st, wh.Verify(req), ErrMessageTimestampTooNew)
	})

	t.Run("KO - verify signature error", func(st *testing.T) {
		id := idx.New("msg")
		ts := fmt.Sprintf("%d", time.Now().UnixMilli())

		req := httptest.NewRequest(http.MethodPost, "/webhook/demo", io.NopCloser(strings.NewReader(body)))
		req.Header.Set(HeaderId, id)
		req.Header.Set(HeaderTimestamp, ts)

		key := idx.Build(DefaultKeyNs, utils.RandomString(128))
		signature := signature.Sign(key, fmt.Sprintf("%s.%s.%s", id, ts, body))
		req.Header.Set(HeaderSignature, signature)

		require.ErrorIs(st, wh.Verify(req), ErrSignatureMismatch)
	})
}

var (
	keys = []string{
		idx.Build(DefaultKeyNs, utils.RandomString(128)),
		idx.Build(DefaultKeyNs, utils.RandomString(128)),
		idx.Build(DefaultKeyNs, utils.RandomString(128)),
	}
)
