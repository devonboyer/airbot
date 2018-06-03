package airbot

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/devonboyer/airbot/botengine"
	"github.com/devonboyer/airbot/messenger"
	log "github.com/sirupsen/logrus"
)

const (
	msgChannelSize    = 1024
	eventsChannelSize = 1024
)

// MessengerService is an adapter for the Messenger API and implements the botengine.ChatService interface.
type MessengerService struct {
	client   *messenger.Client
	msgCh    chan *botengine.Message
	eventsCh chan *messenger.Event

	stopped chan struct{}
	wg      sync.WaitGroup
}

func NewMessengerService(accessToken string) *MessengerService {
	c := messenger.New(
		accessToken,
		messenger.WithHTTPClient(&http.Client{Timeout: 30 * time.Second}),
	)
	return &MessengerService{
		client:   c,
		msgCh:    make(chan *botengine.Message, msgChannelSize),
		eventsCh: make(chan *messenger.Event, eventsChannelSize),
		stopped:  make(chan struct{}),
	}
}

func (svc *MessengerService) Run() {
	for i := 0; i < 10; i++ {
		svc.wg.Add(1)
		go svc.worker()
	}
}

func (svc *MessengerService) worker() {
	defer svc.wg.Done()

	for {
		select {
		case ev := <-svc.eventsCh:
			svc.handleEvent(ev)
		case <-svc.stopped:
			return
		}
	}
}

func (svc *MessengerService) handleEvent(ev *messenger.Event) {
	for _, entry := range ev.Entries {
		callback := entry.Messaging[0]
		if callback.Message != nil {
			log.Infof("Received message (psid: %s, text: %s)", callback.Sender.ID, callback.Message.Text)

			// Mark message as seen.
			ctx := context.Background()
			if err := svc.client.SendByID(callback.Sender.ID).Action(messenger.MarkSeen).Do(ctx); err != nil {
				log.Errorf("Failed to mark seen, %s", err)
			}

			svc.msgCh <- &botengine.Message{
				Sender: botengine.User{ID: callback.Sender.ID},
				Body:   callback.Message.Text,
			}
		}
	}
}

func (svc *MessengerService) Events() chan<- *messenger.Event {
	return svc.eventsCh
}

func (svc *MessengerService) Messages() <-chan *botengine.Message {
	return svc.msgCh
}

func (svc *MessengerService) TypingOn(ctx context.Context, user botengine.User) error {
	return svc.client.
		SendByID(user.ID).
		Action(messenger.TypingOn).
		Do(ctx)
}

func (svc *MessengerService) TypingOff(ctx context.Context, user botengine.User) error {
	return svc.client.
		SendByID(user.ID).
		Action(messenger.TypingOff).
		Do(ctx)
}

func (svc *MessengerService) Send(ctx context.Context, res *botengine.Response) error {
	err := svc.client.
		SendByID(res.Recipient.ID).
		Message().
		Text(res.Body).
		Do(ctx)
	if err != nil {
		log.Infof("Failed to send message (psid: %s, text: %s), %s", res.Recipient.ID, res.Body, err)
		return err
	}
	log.Infof("Sent message (psid: %s, text: %s)", res.Recipient.ID, res.Body)
	return nil
}

func (svc *MessengerService) Close() {
	close(svc.stopped)
	svc.wg.Wait()
}
