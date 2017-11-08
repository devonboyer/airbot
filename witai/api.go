package witai

import (
	"fmt"
	"net/http"
)

const apiVersion = "20170307"

type ClientOption interface {
	Apply(*Client)
}

func WithHTTPClient(client *http.Client) ClientOption {
	return withHTTPClient{client}
}

type withHTTPClient struct{ client *http.Client }

func (w withHTTPClient) Apply(c *Client) {
	c.hc = w.client
}

type Client struct {
	accessToken string
	hc          *http.Client
	basePath    string
}

func New(accessToken string, opts ...ClientOption) *Client {
	o := []ClientOption{
		WithHTTPClient(http.DefaultClient),
	}
	opts = append(o, opts...)
	client := &Client{
		accessToken: accessToken,
		basePath:    fmt.Sprintf("https://api.wit.ai/v=%s", apiVersion),
	}
	for _, opt := range opts {
		opt.Apply(client)
	}
	return client
}

func setAuthorizationHeader(headers http.Header, token string) {
	headers.Set("Authorization", fmt.Sprintf("Bearer %s", token))
}

func setAcceptHeader(headers http.Header, value string) {
	headers.Set("Accept", value)
}

func setContentType(headers http.Header, value string) {
	headers.Set("Content-Type", value)
}
