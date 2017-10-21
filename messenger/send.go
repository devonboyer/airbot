package messenger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	"golang.org/x/net/context/ctxhttp"
)

type NotifType string

const (
	RegularNotif = NotifType("REGULAR")
	SilentNotif  = NotifType("SILENT_PUSH")
	NoNotif      = NotifType("NO_PUSH")
)

type RequestHandle struct {
	client      *Client
	recipientID string
}

func (c *Client) Send(recipientID string) *RequestHandle {
	return &RequestHandle{
		client:      c,
		recipientID: recipientID,
	}
}

type MessageHandle struct {
	client      *Client
	recipientID string
	notifType   NotifType
}

func (r *RequestHandle) Message(notifType NotifType) *MessageHandle {
	return &MessageHandle{
		client:      r.client,
		recipientID: r.recipientID,
		notifType:   notifType,
	}
}

type MessageRequestCall struct {
	client *Client
	req    *MessageRequest
}

func (r *MessageHandle) Text(text string) *MessageRequestCall {
	return &MessageRequestCall{
		client: r.client,
		req: &MessageRequest{
			Recipient: Recipient{ID: r.recipientID},
			Message:   Message{Text: text},
			NotifType: string(r.notifType),
		},
	}
}

func (c *MessageRequestCall) Do(ctx context.Context) error {
	res, err := c.doRequest(ctx)
	if err != nil {
		return err
	}
	return checkResponse(res)
}

func (c *MessageRequestCall) doRequest(ctx context.Context) (*http.Response, error) {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(c.req)
	if err != nil {
		return nil, err
	}
	url := path.Join(c.client.basePath, "messages") + fmt.Sprintf("?accessToken=%s", c.client.accessToken)
	req, _ := http.NewRequest("POST", url, buf)
	setContentType(req.Header, "application/json")
	if ctx == nil {
		return c.client.hc.Do(req)
	}
	return ctxhttp.Do(ctx, c.client.hc, req)
}
