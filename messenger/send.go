package messenger

import "context"

type NotifType string

const (
	RegularNotif = NotifType("REGULAR")
	SilentNotif  = NotifType("SILENT_PUSH")
	NoNotif      = NotifType("NO_PUSH")
)

type RequestHandle struct {
	client      *Client
	recipientID string
	notifType   NotifType
}

func (c *Client) Send(recipientID string) *RequestHandle {
	return &RequestHandle{
		client:      c,
		recipientID: recipientID,
		notifType:   RegularNotif,
	}
}

func (r *RequestHandle) NotifType(notifType NotifType) *RequestHandle {
	r.notifType = notifType
	return r
}

type TextRequestHandle struct {
	*RequestHandle
	text string
}

func (r *RequestHandle) Text(text string) *TextRequestHandle {
	return &TextRequestHandle{
		RequestHandle: r,
		text:          text,
	}
}

func (r *RequestHandle) Do(ctx context.Context) error {
	return nil
}
