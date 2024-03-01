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
		validator.StringAlphaNumericUnderscoreDot("STREAMING.ENTITIES.EVENT.SUBJECT", e.Subject),
		validator.StringRequired("STREAMING.ENTITIES.EVENT.ID", e.Id),
		validator.SliceRequired("STREAMING.ENTITIES.EVENT.DATA", e.Data),
	)
}

func (e *Event) String() string {
	bytes, _ := json.Marshal(e)
	return string(bytes)
}
