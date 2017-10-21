package airtable

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const majorAPIVersion = "0"

type Client struct {
	apiKey   string
	basePath string
	hc       *http.Client
}

func New(apiKey string) *Client {
	return &Client{
		apiKey:   apiKey,
		basePath: fmt.Sprintf("https://api.airtable.com/v%s", majorAPIVersion),
		hc:       http.DefaultClient,
	}
}

var xVersionHeader = fmt.Sprintf("%s.1.0", majorAPIVersion)

func setVersionHeader(headers http.Header) {
	headers.Set("x-api-version", xVersionHeader)
}

type Error struct {
	StatusCode int
	Type       string `json:"type"`
	Message    string `json:"message"`
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
