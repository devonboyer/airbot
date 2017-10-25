package airbot

import (
	"github.com/devonboyer/airbot/botengine"
	"github.com/devonboyer/airbot/messenger"
)

type MessengerSource struct {
	eventsChan chan *botengine.Event
}

func (s *MessengerSource) Events() <-chan *botengine.Event {
	return s.eventsChan
}

// Implements messenger.EventHandler interface
func (s *MessengerSource) HandleEvent(ev *messenger.WebhookEvent) {
	// Convert between event types
}

func (s *MessengerSource) Close() {

}

type MessengerSink struct {
	client *messenger.Client
}

func (s *MessengerSource) Flush(ev *botengine.Event) {

}

func (s *MessengerSource) Close() {

}

// func (m *MessengerSource) Send(reply bot.Reply) {
// 	ctx := context.Background()
// 	m.client.
// 		Send(reply.RecipientID).
// 		Message(messenger.RegularNotif).
// 		Text(reply.Text).
// 		Do(ctx)
// }
