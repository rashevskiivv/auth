package client

import (
	"net/http"
	"time"
)

type Client struct {
	http.Client
}

func NewClient() *Client {
	return &Client{Client: http.Client{Timeout: time.Second * 5}}
}

func (c *Client) Do(req *Request) (*http.Response, error) {
	return c.Client.Do(req.Request)
}
