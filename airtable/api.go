package airtable

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const majorAPIVersion = "0"

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
	apiKey   string
	basePath string
	hc       *http.Client
}

func New(apiKey string, opts ...ClientOption) *Client {
	o := []ClientOption{
		WithHTTPClient(http.DefaultClient),
	}
	opts = append(o, opts...)
	client := &Client{
		apiKey:   apiKey,
		basePath: fmt.Sprintf("https://api.airtable.com/v%s", majorAPIVersion),
	}
	for _, opt := range opts {
		opt.Apply(client)
	}
	return client
}

var xVersionHeader = fmt.Sprintf("%s.1.0", majorAPIVersion)

func setVersionHeader(headers http.Header) {
	headers.Set("x-api-version", xVersionHeader)
}

func setAuthorizationHeader(headers http.Header, token string) {
	headers.Set("Authorization", fmt.Sprintf("Bearer %s", token))
}

type Error struct {
	Type       string `json:"type"`
	Message    string `json:"message"`
	StatusCode int
	Body       string
}

func (e *Error) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s (%d): %s", e.Type, e.StatusCode, e.Message)
	}
	return fmt.Sprintf("Unknown (%d): %s", e.StatusCode, e.Body)
}

type errorReply struct {
	Error *Error `json:"error"`
}

func checkResponse(res *http.Response) error {
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		return nil
	}
	slurp, err := ioutil.ReadAll(res.Body)
	if err == nil {
		jerr := new(errorReply)
		err = json.Unmarshal(slurp, jerr)
		if err == nil && jerr.Error != nil {
			if jerr.Error.StatusCode == 0 {
				jerr.Error.StatusCode = res.StatusCode
			}
			jerr.Error.Body = string(slurp)
			return jerr.Error
		}
	}
	return &Error{
		StatusCode: res.StatusCode,
		Body:       string(slurp),
	}
}
