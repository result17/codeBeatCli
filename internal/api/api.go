package api

import "net/http"

const (
	BaseURL            = ""
	DefaultTimeoutSecs = 60
)

type Client struct {
	baseURL string
	client  *http.Client
}

func NewClient(baseURL string) *Client {
	c := &Client{
		baseURL: baseURL,
		client:  http.DefaultClient,
	}
	return c
}
