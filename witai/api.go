package witai

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		basePath:    "https://api.wit.ai",
	}
	for _, opt := range opts {
		opt.Apply(client)
	}
	return client
}

func setAuthorizationHeader(headers http.Header, token string) {
	headers.Set("Authorization", fmt.Sprintf("Bearer %s", token))
}

func setAcceptHeader(headers http.Header) {
	headers.Set("Accept", fmt.Sprintf("application/vnd.wit.%s+json", apiVersion))
}

type Error struct {
	Message    string
	Code       interface{}
	StatusCode int
	Body       string
}

func (e *Error) Error() string {
	if e.Message != "" && e.Code != nil {
		return fmt.Sprintf("%v (%d): %s", e.Code, e.StatusCode, e.Message)
	} else if e.Message != "" {
		return fmt.Sprintf("Unknown (%d): %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("Unknown (%d): %s", e.StatusCode, e.Body)
}

type errorReply struct {
	Error string      `json:"error"`
	Code  interface{} `json:"code"`
}

func checkResponse(res *http.Response) error {
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		return nil
	}
	slurp, err := ioutil.ReadAll(res.Body)
	outErr := &Error{
		StatusCode: res.StatusCode,
		Body:       string(slurp),
	}
	if err == nil {
		jerr := new(errorReply)
		if json.Unmarshal(slurp, jerr); err == nil {
			outErr.Message = jerr.Error
			outErr.Code = jerr.Code
			return outErr
		}
	}
	return outErr
}
