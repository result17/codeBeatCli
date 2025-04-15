package api

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strings"
)

const (
	BaseURL            = ""
	DefaultTimeoutSecs = 60
)

type Client struct {
	baseURL string
	client  *http.Client

	doFunc func(c *Client, req *http.Request) (*http.Response, error)
}

func NewClient(baseURL string) *Client {
	c := &Client{
		baseURL: baseURL,
		client:  http.DefaultClient,
		doFunc: func(c *Client, req *http.Request) (*http.Response, error) {
			req.Header.Set("Accept", "application/json")
			return c.client.Do(req)
		},
	}
	return c
}

func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	res, err := c.doFunc(c, req)
	if err != nil {
		if strings.HasPrefix(c.baseURL, BaseURL) {
			return nil, err
		}
		var dnsError *net.DNSError
		if !errors.As(err, &dnsError) {
			return nil, err
		}
	}
	return res, nil
}
