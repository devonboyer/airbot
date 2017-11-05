package messenger

import (
	"context"

	"github.com/devonboyer/airbot/botengine"
)

const eventBufferSize = 1024

// Listener implements botengine.Listener interface.
type Listener struct {
	EventHandler

	client     *Client
	eventsChan chan *botengine.Event
}

func NewListener(client *Client) *Listener {
	return &Listener{
		client: client,
	}
}

func (l *Listener) Events() <-chan *botengine.Event {
	return l.eventsChan
}

func (l *Listener) HandleEvent(ev *Event) {
	for _, entry := range ev.Entries {
		callback := entry.Messaging[0]
		if callback.Message != nil {
			l.client.logger.Printf("messenger: Received message (psid: %s, text: %s)", callback.Sender.ID, callback.Message.Text)

			// Mark message as seen.
			ctx := context.Background()
			err := l.client.MarkSeen(ctx, callback.Sender.ID)
			if err != nil {
				l.client.logger.Printf("messenger: Failed to mark seen, %s", err)
			}

			l.eventsChan <- &botengine.Event{
				Type: botengine.MessageEvent,
				Object: &botengine.Message{
					User: botengine.User{ID: callback.Sender.ID},
					Text: callback.Message.Text,
				},
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

func (s *Sender) Send(ctx context.Context, ev *botengine.Event) error {
	switch ev.Type {
	case botengine.MessageEvent:
		msg := ev.Object.(*botengine.Message)
		err := s.client.
			SendByID(msg.User.ID).
			Message(RegularNotif).
			Text(msg.Text).
			Do(ctx)
		if err != nil {
			s.client.logger.Printf("messenger: Failed to send message (psid: %s, text: %s), %s", msg.User.ID, msg.Text, err)
			return err
		}
		s.client.logger.Printf("messenger: Sent message (psid: %s, text: %s)", msg.User.ID, msg.Text)
	default:
	}
	return nil
}

func (s *Sender) Close() {}
