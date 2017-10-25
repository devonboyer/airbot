package botengine

import "testing"

const bufferSize = 1024

type mockSource struct {
	eventsChan chan *Event
}

func newMockSource() *mockSource { return &mockSource{eventsChan: make(chan *Event, bufferSize)} }

func (m *mockSource) Events() <-chan *Event {
	return m.eventsChan
}

func (m *mockSource) Close() error {
	return nil
}

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

}
