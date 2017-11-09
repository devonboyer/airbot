package witai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/net/context/ctxhttp"
)

type MessageCall struct {
	client    *Client
	urlParams url.Values
}

func (c *Client) Message(q string) *MessageCall {
	m := &MessageCall{
		client:    c,
		urlParams: make(url.Values),
	}
	return m.query(q)
}

func (c *MessageCall) query(q string) *MessageCall {
	c.urlParams.Set("q", q)
	return c
}

func (c *MessageCall) Do(ctx context.Context) (*Message, error) {
	v := &Message{}
	res, err := c.doRequest(ctx)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	if err := json.NewDecoder(res.Body).Decode(v); err != nil {
		return nil, err
	}
	return v, nil
}

func (c *MessageCall) doRequest(ctx context.Context) (*http.Response, error) {
	url := fmt.Sprintf("%s/message", c.client.basePath)
	if len(c.urlParams) > 0 {
		url += "?" + c.urlParams.Encode()
	}
	req, _ := http.NewRequest("GET", url, nil)
	setAuthorizationHeader(req.Header, c.client.accessToken)
	setAcceptHeader(req.Header)
	if ctx == nil {
		return c.client.hc.Do(req)
	}
	return ctxhttp.Do(ctx, c.client.hc, req)
}
