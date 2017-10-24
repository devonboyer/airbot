package airbot

import (
	"context"

	"github.com/devonboyer/airbot/botengine"
	"github.com/devonboyer/airbot/messenger"
)

type MessengerQueue struct {
	botengine.Queue
}

func (q *MessengerQueue) Push(ctx context.Context, ev botengine.Event) {
	q.Push(ctx, ev)
}

func (q *MessengerQueue) Pop(ctx context.Context) botengine.Event {
	return q.Pop(ctx)
}

// Implements messenger.EventHandler interface
func (q *MessengerQueue) HandleEvent(ev *messenger.WebhookEvent) {
	// Convert between event types
}

func (q *MessengerQueue) Close() {
	q.Close()
}

// func (m *MessengerSource) Send(reply bot.Reply) {
// 	ctx := context.Background()
// 	m.client.
// 		Send(reply.RecipientID).
// 		Message(messenger.RegularNotif).
// 		Text(reply.Text).
// 		Do(ctx)
// }
