package botengine

import "testing"
import "time"
import "github.com/stretchr/testify/require"
import "io"
import "fmt"

const bufferSize = 1024

type mockSource struct {
	eventsChan chan *Event
}

func newMockSource() *mockSource { return &mockSource{eventsChan: make(chan *Event, bufferSize)} }

func (m *mockSource) Events() <-chan *Event {
	return m.eventsChan
}

func (m *mockSource) Close() {}

type mockSink struct {
	flushedChan chan *Event
}

func newMockSink() *mockSink { return &mockSink{flushedChan: make(chan *Event, bufferSize)} }

func (m *mockSink) Flush(ev *Event) error {
	m.flushedChan <- ev
	return nil
}

func (m *mockSink) Close() {}

func TestEngine(t *testing.T) {
	source := newMockSource()
	sink := newMockSink()
	e := New(source, sink)
	e.Handle("ping", func(w io.Writer, ev *Event) {
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
			"command bot does not undertand",
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
					Text: defaultNotFoundReply,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source.eventsChan <- test.sourceEvent

			select {
			case ev := <-sink.flushedChan:
				require.Equal(t, ev, test.sinkEvent)
			case <-time.After(1 * time.Second):
				require.Fail(t, "timout")
			}
		})
	}
}
