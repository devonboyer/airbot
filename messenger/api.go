package messenger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context/ctxhttp"
)

const apiVersion = "2.8"

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

type Interface interface {
	SendByID(recipientID string) SendInterface
}

type SendInterface interface {
	Action() *SenderActionCall
	Message() *SenderActionCall
}

// There is a client and a server component...maybe separate?

type Client struct {
	accessToken string
	basePath    string
	hc          *http.Client
	skipVerify  bool
}

func New(accessToken string, opts ...ClientOption) *Client {
	o := []ClientOption{
		WithHTTPClient(http.DefaultClient),
	}
	opts = append(o, opts...)
	client := &Client{
		accessToken: accessToken,
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
	url := fmt.Sprintf("%s/messages", c.basePath) + fmt.Sprintf("?access_token=%s", c.accessToken)
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

func handleError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

// https://developers.facebook.com/docs/messenger-platform/reference/send-api/error-codes
type Error struct {
	Message      string `json:"message"`
	Type         string `json:"type"`
	Code         int    `json:"code"`
	ErrorSubcode int    `json:"error_subcode"`
	FBtraceID    string `json:"fbtrace_id"`
	StatusCode   int
	Body         string
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
