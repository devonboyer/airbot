package airbot

import (
	"context"
	"sync"

	"github.com/devonboyer/airbot/bot"
	"github.com/devonboyer/airbot/messenger"
)

type MessengerSource struct {
	client  *messenger.Client
	msgs    chan bot.Message
	stopped chan struct{}
	wg      sync.WaitGroup
}

func NewMessengerSource(client *messenger.Client) *MessengerSource {
	source := &MessengerSource{
		client:  client,
		msgs:    make(chan bot.Message),
		stopped: make(chan struct{}),
		wg:      sync.WaitGroup{},
	}
	go source.convertChannels()
	return source
}

func (m *MessengerSource) convertChannels() {
	m.wg.Add(1)
	defer m.wg.Done()
	for {
		select {
		case msg := <-m.client.Messages():
			m.msgs <- bot.Message{
				SenderID: msg.Sender.ID,
				Text:     msg.Message.Text,
			}
		case <-m.stopped:
			return
		}
	}
}

func (m *MessengerSource) Messages() <-chan bot.Message {
	return m.msgs
}

func (m *MessengerSource) Send(reply bot.Reply) {
	ctx := context.Background()
	m.client.
		Send(reply.RecipientID).
		Message(messenger.RegularNotif).
		Text(reply.Text).
		Do(ctx)
}

func (m *MessengerSource) Stop() {
	close(m.stopped)
	m.wg.Wait()
}
