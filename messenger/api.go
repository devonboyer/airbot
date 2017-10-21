package messenger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const apiVersion = "2.6"

type Client struct {
	accessToken string
	verifyToken string
	appSecret   string
	basePath    string
	hc          *http.Client
}

func New(accessToken, verifyToken, appSecret string) *Client {
	return &Client{
		accessToken: accessToken,
		verifyToken: verifyToken,
		appSecret:   appSecret,
		basePath:    fmt.Sprintf("https://graph.facebook.com/v%s/me", apiVersion),
		hc:          http.DefaultClient,
	}
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
