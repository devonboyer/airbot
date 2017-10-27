package airbot

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/devonboyer/airbot/botengine"
	"github.com/devonboyer/airbot/messenger"
)

const bufferSize = 1024

type MessengerSource struct {
	client     *messenger.Client
	eventsChan chan *botengine.Event
}

func NewMessengerSource(client *messenger.Client) *MessengerSource {
	return &MessengerSource{
		client:     client,
		eventsChan: make(chan *botengine.Event, bufferSize),
	}
}

func (s *MessengerSource) Events() <-chan *botengine.Event {
	return s.eventsChan
}

// HandleEvent implements messenger.EventHandler interface
func (s *MessengerSource) HandleEvent(ev *messenger.WebhookEvent) {
	for _, entry := range ev.Entries {
		callback := entry.Messaging[0]
		if callback.Message != nil {
			logentry := logrus.WithFields(logrus.Fields{
				"psid": callback.Sender.ID,
				"text": callback.Message.Text,
			})

			logentry.Info("received message event")

			ctx := context.Background()
			err := s.client.MarkSeen(ctx, callback.Sender.ID)
			if err != nil {
				logentry.WithError(err).Error("could not mark seen")
			}

			s.eventsChan <- &botengine.Event{
				Type: botengine.MessageEvent,
				Object: &botengine.Message{
					User: botengine.User{ID: callback.Sender.ID},
					Text: callback.Message.Text,
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
	return &MessengerSink{client: client}
}

func (s *MessengerSink) Flush(ev *botengine.Event) error {
	switch ev.Type {
	case botengine.MessageEvent:
		msg := ev.Object.(*botengine.Message)
		logentry := logrus.WithFields(logrus.Fields{
			"psid": msg.User.ID,
			"text": msg.Text,
		})
		ctx := context.Background()
		err := s.client.
			Send(msg.User.ID).
			Message(messenger.RegularNotif).
			Text(msg.Text).
			Do(ctx)
		if err != nil {
			logentry.WithError(err).Error("Failed to send message")
			return err
		}
		logentry.Info("Sent message")
	default:
	}
	return nil
}

func (s *MessengerSink) Close() {}
