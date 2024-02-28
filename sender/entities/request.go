package entities

import "net/http"

type Request struct {
	Method  string      `json:"method"`
	Headers http.Header `json:"headers"`
	Uri     string      `json:"uri"`
	Body    []byte      `json:"body"`
}
