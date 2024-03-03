package entities

import (
	"net/http"

	"github.com/kanthorlabs/common/validator"
)

type Request struct {
	Method  string      `json:"method"`
	Headers http.Header `json:"headers"`
	Uri     string      `json:"uri"`
	Body    []byte      `json:"body"`
}

func (req *Request) Validate() error {
	return validator.Validate(
		validator.StringUri("SENDER.REQUEST.URI", req.Uri),
		validator.StringOneOf("SENDER.REQUEST.METHOD", req.Method, []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch}),
	)
}
