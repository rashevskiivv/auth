package client

import (
	"io"
	"net/http"
)

type Request struct {
	*http.Request
}

func NewRequest(method string, url string, body io.Reader) *Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil
	}
	return &Request{Request: req}
}

func (r *Request) AddAuthHeaders(headers map[string]string) {
	for k, v := range headers {
		if k == "" || v == "" {
			continue
		}
		r.Header.Add(k, v)
	}
}
