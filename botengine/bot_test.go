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
	msgChan chan *Message
}

func newMockListener() *mockListener { return &mockListener{msgChan: make(chan *Message, 1)} }

func (m *mockListener) Messages() <-chan *Message {
	return m.msgChan
}

func (m *mockListener) Close() {}

type mockSender struct {
	sentChan chan *Response
}

func newMockSender() *mockSender { return &mockSender{sentChan: make(chan *Response, 1)} }

func (m *mockSender) TypingOn(_ context.Context, _ User) error { return nil }

func (m *mockSender) TypingOff(_ context.Context, _ User) error { return nil }

func (m *mockSender) Send(_ context.Context, res *Response) error {
	m.sentChan <- res
	return nil
}

func (m *mockSender) Close() {}

func Test_Engine(t *testing.T) {
	listener := newMockListener()
	sender := newMockSender()

	e := New()
	e.Listener = listener
	e.Sender = sender
	e.NotFoundHandler = HandlerFunc(func(w io.Writer, req *Request) {
		fmt.Fprintf(w, req.Message.Body)
	})

	e.HandleFunc("ping", func(w io.Writer, _ *Request) {
		fmt.Fprintf(w, "pong")
	})
	e.Run()
	defer e.Stop()

	tests := []struct {
		name string
		msg  *Message
		res  *Response
	}{
		{
			"command bot understands",
			&Message{
				Sender: User{ID: "1"},
				Body:   "ping",
			},
			&Response{
				Recipient: User{ID: "1"},
				Body:      "pong",
			},
		},
		{
			"not found",
			&Message{
				Sender: User{ID: "1"},
				Body:   "foo",
			},
			&Response{
				Recipient: User{ID: "1"},
				Body:      "foo",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			listener.msgChan <- test.msg
			select {
			case res := <-sender.sentChan:
				require.Equal(t, res, test.res)
			case <-time.After(1 * time.Second):
				require.Fail(t, "timout")
			}
		})
	}
}
