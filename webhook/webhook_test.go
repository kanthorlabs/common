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

	t.Run("KO - no keys error", func(st *testing.T) {
		_, err := New([]string{})
		require.ErrorContains(st, err, "must not be empty")
	})

	t.Run("KO - too many keys error", func(st *testing.T) {
		keys := make([]string, MaxKeys+1)
		for i := range keys {
			keys[i] = idx.Build(IdNsEpSec, utils.RandomString(128))
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

	t.Run("KO - verify timestamp error", func(st *testing.T) {
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

	t.Run("KO - verify signature error", func(st *testing.T) {
		id := idx.New("msg")
		ts := fmt.Sprintf("%d", time.Now().UnixMilli())

		req := httptest.NewRequest(http.MethodPost, "/webhook/demo", io.NopCloser(strings.NewReader(body)))
		req.Header.Set(HeaderId, id)
		req.Header.Set(HeaderTimestamp, ts)

		key := idx.Build(IdNsEpSec, utils.RandomString(128))
		signature := signature.Sign(key, fmt.Sprintf("%s.%s.%s", id, ts, body))
		req.Header.Set(HeaderSignature, signature)

		require.ErrorIs(st, wh.Verify(req), ErrSignatureMismatch)
	})
}

func TestVerifyTimestamp(t *testing.T) {
	wh, err := New(keys)
	require.NoError(t, err)

	t.Run("OK", func(st *testing.T) {
		ts := fmt.Sprintf("%d", time.Now().UnixMilli())
		options := &VerifyOptions{TimestampToleranceIgnore: true}
		err = wh.VerifyTimestamp(ts, options)
		require.NoError(st, err)
	})

	t.Run("OK - ignore", func(st *testing.T) {
		ts := fmt.Sprintf("%d", time.Now().Add(-time.Hour*365).UnixMilli())
		options := &VerifyOptions{TimestampToleranceIgnore: true}
		err = wh.VerifyTimestamp(ts, options)
		require.NoError(st, err)
	})

	t.Run("KO - parse timestamp error", func(st *testing.T) {
		ts := fmt.Sprintf("xxx-%d", time.Now().UnixMilli())
		options := &VerifyOptions{TimestampToleranceIgnore: true}
		err = wh.VerifyTimestamp(ts, options)
		require.ErrorIs(st, err, ErrMessageTimestampMalformed)
	})

	t.Run("OK - too old error", func(st *testing.T) {
		ts := fmt.Sprintf("%d", time.Now().Add(-DefaultToleranceDuration-time.Hour).UnixMilli())
		options := &VerifyOptions{TimestampToleranceDuration: DefaultToleranceDuration}
		err = wh.VerifyTimestamp(ts, options)
		require.ErrorIs(st, err, ErrMessageTimestampTooOld)
	})

	t.Run("OK - too new error", func(st *testing.T) {
		ts := fmt.Sprintf("%d", time.Now().Add(DefaultToleranceDuration+time.Hour).UnixMilli())
		options := &VerifyOptions{TimestampToleranceDuration: DefaultToleranceDuration}
		err = wh.VerifyTimestamp(ts, options)
		require.ErrorIs(st, err, ErrMessageTimestampTooNew)
	})
}

func TestVerifySignature(t *testing.T) {
	wh, err := New(keys)
	require.NoError(t, err)

	id := idx.New("msg")
	ts := fmt.Sprintf("%d", time.Now().UnixMilli())
	body := utils.Stringify(map[string]any{
		"app_id": "msg_2ePVr2tTfiJA20mN8wkc8EkGZu4",
		"type":   "testing.openapi",
		"object": map[string]any{"from_client": "openapi", "say": "hello"},
	})

	t.Run("OK", func(st *testing.T) {
		key := testdata.Fake.RandomStringElement(keys)
		expected := signature.Sign(key, fmt.Sprintf("%s.%s.%s", id, ts, body))

		err := wh.VerifySignature(id, ts, body, expected)
		require.NoError(st, err)
	})

	t.Run("KO", func(st *testing.T) {
		key := idx.Build(IdNsEpSec, utils.RandomString(128))
		expected := signature.Sign(key, fmt.Sprintf("%s.%s.%s", id, ts, body))

		err := wh.VerifySignature(id, ts, body, expected)
		require.ErrorIs(st, err, ErrSignatureMismatch)
	})
}

var (
	keys = []string{
		idx.Build(IdNsEpSec, utils.RandomString(128)),
		idx.Build(IdNsEpSec, utils.RandomString(128)),
		idx.Build(IdNsEpSec, utils.RandomString(128)),
	}
)
