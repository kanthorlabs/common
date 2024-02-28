package entities

import "net/http"

type Response struct {
	Status  int         `json:"status"`
	Headers http.Header `json:"headers"`
	Uri     string      `json:"uri"`
	Body    []byte      `json:"body"`
}
