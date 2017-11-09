package botengine

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type mockChatService struct {
	msgChan  chan *Message
	sentChan chan *Response
}

func newMockChatService() *mockChatService {
	return &mockChatService{
		msgChan:  make(chan *Message, 1),
		sentChan: make(chan *Response, 1),
	}
}

func (m *mockChatService) Messages() <-chan *Message {
	return m.msgChan
}

func (m *mockChatService) TypingOn(_ context.Context, _ User) error { return nil }

func (m *mockChatService) TypingOff(_ context.Context, _ User) error { return nil }

func (m *mockChatService) Send(_ context.Context, res *Response) error {
	m.sentChan <- res
	return nil
}

func (m *mockChatService) Close() {}

func Test_Engine(t *testing.T) {
	chatService := newMockChatService()

	e := New()
	e.ChatService = chatService
	e.NotFoundHandler = HandlerFunc(func(w io.Writer, msg *Message) {
		fmt.Fprintf(w, msg.Body)
	})

	e.HandleFunc("ping", func(w io.Writer, _ *Message) {
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
			chatService.msgChan <- test.msg
			select {
			case res := <-chatService.sentChan:
				require.Equal(t, res, test.res)
			case <-time.After(1 * time.Second):
				require.Fail(t, "timout")
			}
		})
	}
}
