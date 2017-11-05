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

func (c *Client) SendByID(recipientID string) *SendHandle {
	return &SendHandle{
		client:      c,
		recipientID: recipientID,
	}
}

func (c *Client) MarkSeen(ctx context.Context, recipientID string) error {
	return c.SendByID(recipientID).Action(MarkSeen).Do(ctx)
}

func (c *Client) TypingOn(ctx context.Context, recipientID string) error {
	return c.SendByID(recipientID).Action(TypingOn).Do(ctx)
}

func (c *Client) TypingOff(ctx context.Context, recipientID string) error {
	return c.SendByID(recipientID).Action(TypingOff).Do(ctx)
}

type SendHandle struct {
	client      *Client
	recipientID string
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

func (r *SendHandle) Message(notifType NotifType) *MessageHandle {
	return &MessageHandle{
		client:      r.client,
		recipientID: r.recipientID,
		notifType:   notifType,
	}
}

type SenderActionCall struct {
	client *Client
	data   *SenderActionMarshaler
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

type SendMessageCall struct {
	client *Client
	data   *MessageMarshaler
}

func (c *SendMessageCall) Do(ctx context.Context) error {
	res, err := c.client.doRequest(ctx, c.data)
	if err != nil {
		return err
	}
	return checkResponse(res)
}
