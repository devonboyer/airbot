package airbot

import (
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

// Implements messenger.EventHandler interface
func (s *MessengerSource) HandleEvent(ev *messenger.WebhookEvent) {
	// Convert between event types
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
	// ctx := context.Background()
	// s.client.
	// 	Send(reply.RecipientID).
	// 	Message(messenger.RegularNotif).
	// 	Text(reply.Text).
	// 	Do(ctx)
	return nil
}

func (s *MessengerSink) Close() {}
