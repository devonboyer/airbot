package messenger

import "context"

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
	client    *Client
	recipient Recipient
}

func (c *Client) Send(recipientID string) *SendHandle {
	return &SendHandle{
		client:    c,
		recipient: Recipient{ID: recipientID},
	}
}

type SenderActionCall struct {
	client *Client
	body   *SenderActionBody
}

func (r *SendHandle) Action(action SenderAction) *SenderActionCall {
	return &SenderActionCall{
		client: r.client,
		body: &SenderActionBody{
			Recipient: r.recipient,
			Action:    string(action),
		},
	}
}

func (c *SenderActionCall) Do(ctx context.Context) error {
	res, err := c.client.doRequest(ctx, c.body)
	if err != nil {
		return err
	}
	return checkResponse(res)
}

type MessageHandle struct {
	client    *Client
	recipient Recipient
	notifType NotifType
}

func (r *SendHandle) Message(notifType NotifType) *MessageHandle {
	return &MessageHandle{
		client:    r.client,
		recipient: r.recipient,
		notifType: notifType,
	}
}

type SendMessageCall struct {
	client *Client
	body   *SendBody
}

func (r *MessageHandle) Text(text string) *SendMessageCall {
	return &SendMessageCall{
		client: r.client,
		body: &SendBody{
			Recipient: r.recipient,
			Message:   Message{Text: text},
			NotifType: string(r.notifType),
		},
	}
}

func (c *SendMessageCall) Do(ctx context.Context) error {
	res, err := c.client.doRequest(ctx, c.body)
	if err != nil {
		return err
	}
	return checkResponse(res)
}
