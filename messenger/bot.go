package messenger

import (
	"context"

	"github.com/devonboyer/airbot/botengine"
)

const msgBufferSize = 1024

// Listener implements botengine.Listener and EventHandler interfaces.
type Listener struct {
	EventHandler

	client  *Client
	msgChan chan *botengine.Message
}

func NewListener(client *Client) *Listener {
	return &Listener{
		client:  client,
		msgChan: make(chan *botengine.Message, msgBufferSize),
	}
}

func (l *Listener) Messages() <-chan *botengine.Message {
	return l.msgChan
}

func (l *Listener) HandleEvent(ev *Event) {
	for _, entry := range ev.Entries {
		callback := entry.Messaging[0]
		if callback.Message != nil {
			l.client.logger.Printf("messenger: Received message (psid: %s, text: %s)", callback.Sender.ID, callback.Message.Text)

			// Mark message as seen.
			ctx := context.Background()
			if err := l.client.SendByID(callback.Sender.ID).Action(MarkSeen).Do(ctx); err != nil {
				l.client.logger.Printf("messenger: Failed to mark seen, %s", err)
			}

			l.msgChan <- &botengine.Message{
				Sender: botengine.User{ID: callback.Sender.ID},
				Body:   callback.Message.Text,
			}
		}
	}
}

func (l *Listener) Close() {}

// Sender implements botengine.Sender interface.
type Sender struct {
	client *Client
}

func NewSender(client *Client) *Sender {
	return &Sender{
		client: client,
	}
}

func (s *Sender) TypingOn(ctx context.Context, user botengine.User) error {
	return s.client.SendByID(user.ID).Action(TypingOn).Do(ctx)
}

func (s *Sender) TypingOff(ctx context.Context, user botengine.User) error {
	return s.client.SendByID(user.ID).Action(TypingOff).Do(ctx)
}

func (s *Sender) Send(ctx context.Context, res *botengine.Response) error {
	err := s.client.
		SendByID(res.Recipient.ID).
		Message(RegularNotif).
		Text(res.Body).
		Do(ctx)
	if err != nil {
		s.client.logger.Printf("messenger: Failed to send message (psid: %s, text: %s), %s", res.Recipient.ID, res.Body, err)
		return err
	}
	s.client.logger.Printf("messenger: Sent message (psid: %s, text: %s)", res.Recipient.ID, res.Body)
	return nil
}

func (s *Sender) Close() {}
