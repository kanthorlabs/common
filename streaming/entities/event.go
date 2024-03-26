package entities

import (
	"encoding/json"

	"github.com/kanthorlabs/common/validator"
)

var (
	MetaTrace = "Traceparent"
)

type Event struct {
	Subject string `json:"subject"`

	Id       string            `json:"id"`
	Data     []byte            `json:"data"`
	Metadata map[string]string `json:"metadata"`
}

func (e *Event) Validate() error {
	return validator.Validate(
		validator.StringAlphaNumericUnderscoreHyphenDot("STREAMING.EVENT.SUBJECT", e.Subject),
		validator.StringRequired("STREAMING.EVENT.ID", e.Id),
		validator.SliceRequired("STREAMING.EVENT.DATA", e.Data),
	)
}

func (e *Event) String() string {
	bytes, _ := json.Marshal(e)
	return string(bytes)
}
