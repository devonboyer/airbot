package messenger

type MessageHandle struct {
	client    *Client
	recipient string
	text      string
}

func (c *Client) Message() *MessageHandle {
	return &MessageHandle{
		client: c,
	}
}

func (c *Client) Template() *MessageHandle {}

// attachment and payload diff for each template

// Message,Button,Receipt..all different payloads sort of

type MessageCall struct {
	// fields
}

// recipient, message

// Many different kinds of templates
//
