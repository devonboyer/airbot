package botengine

import (
	"context"
	"fmt"
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

func Test_Bot(t *testing.T) {
	chatService := newMockChatService()

	bot := New()
	bot.ChatService = chatService
	bot.NotFoundHandler = HandlerFunc(func(w ResponseWriter, msg *Message) {
		fmt.Fprintf(w, msg.Body)
		w.SetStatus(StatusNotFound)
	})

	HandleFunc("ping", func(w ResponseWriter, _ *Message) {
		fmt.Fprintf(w, "pong")
	})
	bot.Run()
	defer bot.Stop()

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
				Status:    StatusOk,
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
				Status:    StatusNotFound,
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
