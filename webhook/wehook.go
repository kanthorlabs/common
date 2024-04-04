package webhook

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/kanthorlabs/common/cipher/signature"
	"github.com/kanthorlabs/common/validator"
)

func New(keys []string, withOptions ...Option) (Webhook, error) {
	options := &Options{
		KeyNamespace: DefaultKeyNs,
	}
	for i := range withOptions {
		withOptions[i](options)
	}

	err := validator.Validate(
		validator.SliceRequired("keys", keys),
		validator.SliceMaxLength("keys", keys, MaxKeys),
		validator.Slice(keys, func(i int, item *string) error {
			return validator.StringStartsWith(fmt.Sprintf("keys[%d]", i), *item, options.KeyNamespace)()
		}),
	)
	if err != nil {
		return nil, err
	}

	return &webhook{keys: keys}, nil
}

type Webhook interface {
	Sign(id, ts, body string) []string
	Verify(req *http.Request, withOptions ...VerifyOption) error
}

type webhook struct {
	keys []string
}

func (wh *webhook) Sign(id, ts, body string) []string {
	var signatures []string
	for i := range wh.keys {
		signatures = append(signatures, signature.Sign(wh.keys[i], fmt.Sprintf("%s.%s.%s", id, ts, body)))
	}
	return signatures
}

func (wh *webhook) Verify(req *http.Request, withOptions ...VerifyOption) error {
	options := &VerifyOptions{
		TimestampToleranceDuration: ToleranceDurationDefault,
	}
	for i := range withOptions {
		withOptions[i](options)
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	timestamp := req.Header.Get(HeaderTimestamp)
	if err := wh.verifyTimestamp(timestamp, options); err != nil {
		return err
	}

	id := req.Header.Get(HeaderId)
	expected := req.Header.Get(HeaderSignature)

	data := fmt.Sprintf("%s.%s.%s", id, timestamp, body)
	if err := signature.VerifyAny(wh.keys, data, expected); err != nil {
		return ErrSignatureMismatch
	}

	return nil
}

func (wh *webhook) verifyTimestamp(ts string, options *VerifyOptions) error {
	timestamp, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return ErrMessageTimestampMalformed
	}

	if options.TimestampToleranceIgnore {
		return nil
	}

	t := time.UnixMilli(timestamp)
	low := t.Add(-options.TimestampToleranceDuration).UnixMilli()
	high := t.Add(options.TimestampToleranceDuration).UnixMilli()
	now := time.Now().UnixMilli()

	if now > high {
		return ErrMessageTimestampTooOld
	}
	if now < low {
		return ErrMessageTimestampTooNew
	}

	return nil
}
