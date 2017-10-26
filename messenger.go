package airbot

import (
	"context"

	"github.com/devonboyer/airbot/botengine"
	"github.com/devonboyer/airbot/messenger"
)

const bufferSize = 1024

type MessengerSource struct {
	eventsChan chan *botengine.Event
}

func NewMessengerSource() *MessengerSource {
	return &MessengerSource{
		eventsChan: make(chan *botengine.Event, bufferSize),
	}
}

func (s *MessengerSource) Events() <-chan *botengine.Event {
	return s.eventsChan
}

// HandleEvent implements messenger.EventHandler interface
func (s *MessengerSource) HandleEvent(ev *messenger.WebhookEvent) {
	for _, entry := range ev.Entries {
		switch msg := entry.Messaging[0].(type) {
		case *messenger.MessageEvent:
			s.eventsChan <- &botengine.Event{
				Type: botengine.MessageEvent,
				Object: &botengine.Message{
					User: botengine.User{ID: msg.Sender.ID},
					Text: msg.Message.Text,
				},
			}
		}
	}
}

func (s *MessengerSource) Close() {}

type MessengerSink struct {
	client *messenger.Client
}

func NewMessengerSink(client *messenger.Client) *MessengerSink {
	return &MessengerSink{
		client: client,
	}
}

func (s *MessengerSink) Flush(ev *botengine.Event) error {
	switch ev.Type {
	case botengine.MessageEvent:
		msg := ev.Object.(*botengine.Message)
		ctx := context.Background()
		return s.client.
			Send(msg.User.ID).
			Message(messenger.RegularNotif).
			Text(msg.Text).
			Do(ctx)
	default:
	}
	return nil
}

func (s *MessengerSink) Close() {}
