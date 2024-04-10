package entities

import "net/http"

type Response struct {
	Status  int         `json:"status"`
	Headers http.Header `json:"headers"`
	Uri     string      `json:"uri"`
	Body    []byte      `json:"body"`
}

func (entity *Response) Ok() bool {
	if entity.Status < http.StatusOK {
		return false
	}
	if entity.Status >= http.StatusBadRequest {
		return false
	}

	return true
}

func (entity *Response) StatusText() string {
	if entity.Status == -1 {
		return string(entity.Body)
	}

	return http.StatusText(entity.Status)
}
