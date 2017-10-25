package messenger

import (
	"context"
)

type (
	NotifType    string
	SenderAction string
)

const (
	RegularNotif = NotifType("REGULAR")
	SilentNotif  = NotifType("SILENT_PUSH")
	NoNotif      = NotifType("NO_PUSH")
	TypingOn     = SenderAction("typing_on")
	TypingOff    = SenderAction("typing_off")
	MarkSeen     = SenderAction("mark_seen")
)

type SendHandle struct {
	client      *Client
	recipientID string
}

func (c *Client) Send(recipientID string) *SendHandle {
	return &SendHandle{
		client:      c,
		recipientID: recipientID,
	}
}

type SenderActionCall struct {
	client *Client
	data   *SenderActionMarshaler
}

func (r *SendHandle) Action(action SenderAction) *SenderActionCall {
	return &SenderActionCall{
		client: r.client,
		data: &SenderActionMarshaler{
			Recipient: Recipient{ID: r.recipientID},
			Action:    string(action),
		},
	}
}

func (c *SenderActionCall) Do(ctx context.Context) error {
	res, err := c.client.doRequest(ctx, c.data)
	if err != nil {
		return err
	}
	return checkResponse(res)
}

type MessageHandle struct {
	client      *Client
	recipientID string
	notifType   NotifType
}

func (r *SendHandle) Message(notifType NotifType) *MessageHandle {
	return &MessageHandle{
		client:      r.client,
		recipientID: r.recipientID,
		notifType:   notifType,
	}
}

type SendMessageCall struct {
	client *Client
	data   *MessageMarshaler
}

func (r *MessageHandle) Text(text string) *SendMessageCall {
	return &SendMessageCall{
		client: r.client,
		data: &MessageMarshaler{
			Recipient: Recipient{ID: r.recipientID},
			Message:   Message{Text: text},
			NotifType: string(r.notifType),
		},
	}
}

func (c *SendMessageCall) Do(ctx context.Context) error {
	res, err := c.client.doRequest(ctx, c.data)
	if err != nil {
		return err
	}
	return checkResponse(res)
}
