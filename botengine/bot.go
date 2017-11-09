package botengine

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

// Message represents an incoming message.
type Message struct {
	// Sender is the user who sent the message.
	Sender User
	// Body is the body of the message.
	Body string
}

type Response struct {
	// Recipient is the user who should receive the message.
	Recipient User
	// Body is the body of the message.
	Body string
}

// User is a user in the chat service.
type User struct {
	// ID is the user's unique ID.
	ID string
}

type Handler interface {
	Handle(io.Writer, *Message)
}

type HandlerFunc func(io.Writer, *Message)

func (f HandlerFunc) Handle(w io.Writer, msg *Message) {
	f(w, msg)
}

func NotFound(w io.Writer, msg *Message) {
	fmt.Fprintf(w, fmt.Sprintf("I don't understand \"%s\".", msg.Body))
}

func NotFoundHandler() Handler { return HandlerFunc(NotFound) }

// ChatService is an interface for sending and receiving messages from a chat service
// like Messenger.
//
// A ChatService must be safe for concurrent use by multiple
// goroutines.
type ChatService interface {
	Messages() <-chan *Message
	Send(context.Context, *Response) error
	TypingOn(context.Context, User) error
	TypingOff(context.Context, User) error
	Close()
}

type matcher interface {
	MatchString(string) bool
}

type stringMatcher string

func (m stringMatcher) MatchString(s string) bool {
	return strings.ToLower(s) == string(m)
}

type handlerEntry struct {
	matcher matcher
	handler Handler
}

// Bot receives events from a Listener, dispatches events to handlers, and sends
// responses back to a Sender.
type Bot struct {
	ChatService   ChatService
	NumGoroutines int
	// NotFoundHandler will be called when no handlers match an incoming message.
	NotFoundHandler Handler

	mu       sync.Mutex
	handlers []*handlerEntry
	stopped  chan struct{}
	wg       sync.WaitGroup
}

func New() *Bot {
	return &Bot{
		NumGoroutines:   1,
		NotFoundHandler: NotFoundHandler(),
		mu:              sync.Mutex{},
		handlers:        make([]*handlerEntry, 0),
		stopped:         make(chan struct{}),
		wg:              sync.WaitGroup{},
	}
}

func (b *Bot) Handle(pattern string, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers = append(b.handlers, &handlerEntry{
		matcher: stringMatcher(pattern),
		handler: handler,
	})
}

func (b *Bot) HandleFunc(pattern string, handler func(io.Writer, *Message)) {
	b.Handle(pattern, HandlerFunc(handler))
}

func (b *Bot) Run() <-chan error {
	if b.ChatService == nil {
		panic("botengine: ChatService must not be nil")
	}

	outError := make(chan error, 1)
	for i := 0; i < b.NumGoroutines; i++ {
		b.wg.Add(1)
		go b.run(outError)
	}
	return outError
}

func (b *Bot) run(outError chan error) {
	defer b.wg.Done()

	for {
		select {
		case msg := <-b.ChatService.Messages():
			b.receive(outError, msg)
		case <-b.stopped:
			return
		}
	}
}

func (b *Bot) receive(outError chan error, msg *Message) {
	for _, h := range b.handlers {
		if h.matcher.MatchString(msg.Body) {
			b.dispatch(outError, h.handler, msg)
			return
		}
	}
	if b.NotFoundHandler != nil {
		b.dispatch(outError, b.NotFoundHandler, msg)
	}
}

func (b *Bot) dispatch(outError chan error, handler Handler, msg *Message) {
	buf := &bytes.Buffer{}

	// Call handler.
	ctx := context.Background()
	if err := b.ChatService.TypingOn(ctx, msg.Sender); err != nil {
		outError <- errors.Wrap(err, "could not send typing on action")
	}
	handler.Handle(buf, msg)
	if err := b.ChatService.TypingOff(ctx, msg.Sender); err != nil {
		outError <- errors.Wrap(err, "could not send typing off action")
	}
	if body := buf.String(); body != "" {
		b.send(outError, msg.Sender, body)
	}
}

func (b *Bot) send(outError chan error, recipient User, body string) {
	res := &Response{
		Recipient: recipient,
		Body:      body,
	}
	ctx := context.Background()
	if err := b.ChatService.Send(ctx, res); err != nil {
		outError <- errors.Wrap(err, "could not send response")
	}
}

func (b *Bot) Stop() {
	close(b.stopped)
	b.wg.Wait()

	b.ChatService.Close()
}
