package airbot

import (
	"context"
	"sync"

	"github.com/devonboyer/airbot/bot"
	"github.com/devonboyer/airbot/messenger"
)

type MessengerSender struct {
	client *messenger.Client
}

func NewMessengerSender(client *messenger.Client) *MessengerSender {
	return &MessengerSender{
		client: client,
	}
}

func (m *MessengerSender) Send(reply bot.Reply) {
	ctx := context.Background()
	m.client.
		Send(reply.RecipientID).
		Message(messenger.RegularNotif).
		Text(reply.Text).
		Do(ctx)
}

type MessengerListener struct {
	client  *messenger.Client
	msgs    chan bot.Message
	stopped chan struct{}
	wg      sync.WaitGroup
}

func NewMessengerListener(client *messenger.Client) *MessengerListener {
	listener := &MessengerListener{
		client:  client,
		msgs:    make(chan bot.Message),
		stopped: make(chan struct{}),
		wg:      sync.WaitGroup{},
	}
	go listener.convertChannels()
	return listener
}

func (m *MessengerListener) convertChannels() {
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

func (m *MessengerListener) Messages() <-chan bot.Message {
	return m.msgs
}

func (m *MessengerListener) Stop() {
	close(m.stopped)
	m.wg.Wait()
}
