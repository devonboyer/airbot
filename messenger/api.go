package messenger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"

	"golang.org/x/net/context/ctxhttp"
)

const apiVersion = "2.6"

type logger interface {
	Printf(string, ...interface{})
}

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

func WithLogger(logger logger) ClientOption {
	return withLogger{logger}
}

type withLogger struct{ logger logger }

func (w withLogger) Apply(c *Client) {
	c.logger = w.logger
}

var nopLogger = log.New(ioutil.Discard, "", 0)

type Client struct {
	accessToken string
	verifyToken string
	appSecret   string
	basePath    string
	hc          *http.Client
	logger      logger
	skipVerify  bool
}

func New(accessToken, verifyToken, appSecret string, opts ...ClientOption) *Client {
	o := []ClientOption{
		WithHTTPClient(http.DefaultClient),
		WithLogger(nopLogger),
	}
	opts = append(o, opts...)
	client := &Client{
		accessToken: accessToken,
		verifyToken: verifyToken,
		appSecret:   appSecret,
		basePath:    fmt.Sprintf("https://graph.facebook.com/v%s/me", apiVersion),
	}
	for _, opt := range opts {
		opt.Apply(client)
	}
	return client
}

func (c *Client) doRequest(ctx context.Context, v interface{}) (*http.Response, error) {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(v)
	if err != nil {
		return nil, err
	}
	url := path.Join(c.basePath, "messages") + fmt.Sprintf("?accessToken=%s", c.accessToken)
	req, _ := http.NewRequest("POST", url, buf)
	setContentType(req.Header, "application/json")
	if ctx == nil {
		return c.hc.Do(req)
	}
	return ctxhttp.Do(ctx, c.hc, req)
}

func setContentType(headers http.Header, value string) {
	headers.Set("Content-Type", value)
}

// https://developers.facebook.com/docs/messenger-platform/reference/send-api/error-codes
type Error struct {
	Message    string `json:"message"`
	Type       string `json:"type"`
	Code       int    `json:"code"`
	ErrorData  string `json:"error_data"`
	FBstraceID string `json:"fbstrace_id"`
	StatusCode int
	Body       string
}

func (e *Error) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s (%d): %s", e.Type, e.Code, e.Message)
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
