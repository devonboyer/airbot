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

type SendHandle struct {
	client      *Client
	recipientID string
}

func (h *SendHandle) Action(action SenderAction) *SenderActionCall {
	return &SenderActionCall{
		client: h.client,
		req: &SenderActionRequest{
			Recipient: Recipient{ID: h.recipientID},
			Action:    string(action),
		},
	}
}

func (h *SendHandle) Message() *MessageHandle {
	return &MessageHandle{
		client:      h.client,
		recipientID: h.recipientID,
	}
}

type SenderActionCall struct {
	client *Client
	req    *SenderActionRequest
}

func (c *SenderActionCall) Do(ctx context.Context) error {
	res, err := c.client.doRequest(ctx, c.req)
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

func (h *MessageHandle) NotifType(notifType NotifType) *MessageHandle {
	h.notifType = notifType
	return h
}

func (h *MessageHandle) Text(text string) *SendMessageCall {
	return &SendMessageCall{
		client: h.client,
		req: &MessageRequest{
			Recipient: Recipient{ID: h.recipientID},
			Message:   Message{Text: text},
			NotifType: string(h.notifType),
		},
	}
}

type SendMessageCall struct {
	client *Client
	req    *MessageRequest
}

func (c *SendMessageCall) Do(ctx context.Context) error {
	res, err := c.client.doRequest(ctx, c.req)
	if err != nil {
		return err
	}
	return checkResponse(res)
}
