package botengine

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type mockListener struct {
	eventsChan chan *Event
}

func newMockListener() *mockListener { return &mockListener{eventsChan: make(chan *Event, 1)} }

func (m *mockListener) Events() <-chan *Event {
	return m.eventsChan
}

func (m *mockListener) Close() {}

type mockSender struct {
	sentChan chan *Event
}

func newMockSender() *mockSender { return &mockSender{sentChan: make(chan *Event, 1)} }

func (m *mockSender) Send(_ context.Context, ev *Event) error {
	m.sentChan <- ev
	return nil
}

func (m *mockSender) Close() {}

func Test_Engine(t *testing.T) {
	listener := newMockListener()
	sender := newMockSender()
	e := New(listener, sender, DefaultSettings)

	e.HandleFunc("ping", func(w io.Writer, ev *Event) {
		fmt.Fprintf(w, "pong")
	})
	e.Run()
	defer e.Stop()

	tests := []struct {
		name                   string
		sourceEvent, sinkEvent *Event
	}{
		{
			"command bot understands",
			&Event{
				Type: MessageEvent,
				Object: &Message{
					User: User{ID: "1"},
					Text: "ping",
				},
			},
			&Event{
				Type: MessageEvent,
				Object: &Message{
					User: User{ID: "1"},
					Text: "pong",
				},
			},
		},
		{
			"not found",
			&Event{
				Type: MessageEvent,
				Object: &Message{
					User: User{ID: "1"},
					Text: "foo",
				},
			},
			&Event{
				Type: MessageEvent,
				Object: &Message{
					User: User{ID: "1"},
					Text: DefaultSettings.NotFoundReply,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			listener.eventsChan <- test.sourceEvent
			select {
			case ev := <-sender.sentChan:
				require.Equal(t, ev, test.sinkEvent)
			case <-time.After(1 * time.Second):
				require.Fail(t, "timout")
			}
		})
	}
}
