package messenger

import (
	"context"

	"github.com/devonboyer/airbot/botengine"
)

const msgBufferSize = 1024

// ChatService implements botengine.ChatService and EventHandler interfaces.
type ChatService struct {
	EventHandler

	client  *Client
	msgChan chan *botengine.Message
}

func NewChatService(client *Client) *ChatService {
	return &ChatService{
		client:  client,
		msgChan: make(chan *botengine.Message, msgBufferSize),
	}
}

func (c *ChatService) Messages() <-chan *botengine.Message {
	return c.msgChan
}

func (c *ChatService) HandleEvent(ev *Event) {
	for _, entry := range ev.Entries {
		callback := entry.Messaging[0]
		if callback.Message != nil {
			c.client.logger.Printf("messenger: Received message (psid: %s, text: %s)", callback.Sender.ID, callback.Message.Text)

			// Mark message as seen.
			ctx := context.Background()
			if err := c.client.SendByID(callback.Sender.ID).Action(MarkSeen).Do(ctx); err != nil {
				c.client.logger.Printf("messenger: Failed to mark seen, %s", err)
			}

			c.msgChan <- &botengine.Message{
				Sender: botengine.User{ID: callback.Sender.ID},
				Body:   callback.Message.Text,
			}
		}
	}
}

func (c *ChatService) TypingOn(ctx context.Context, user botengine.User) error {
	return c.client.
		SendByID(user.ID).
		Action(TypingOn).
		Do(ctx)
}

func (c *ChatService) TypingOff(ctx context.Context, user botengine.User) error {
	return c.client.
		SendByID(user.ID).
		Action(TypingOff).
		Do(ctx)
}

func (c *ChatService) Send(ctx context.Context, res *botengine.Response) error {
	err := c.client.
		SendByID(res.Recipient.ID).
		Message().
		Text(res.Body).
		Do(ctx)
	if err != nil {
		c.client.logger.Printf("messenger: Failed to send message (psid: %s, text: %s), %s", res.Recipient.ID, res.Body, err)
		return err
	}
	c.client.logger.Printf("messenger: Sent message (psid: %s, text: %s)", res.Recipient.ID, res.Body)
	return nil
}

func (c *ChatService) Close() {}
