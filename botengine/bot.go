package botengine

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"regexp"
	"sync"
)

type Handler interface {
	Handle(io.Writer, *Message)
}

type HandlerFunc func(io.Writer, *Message)

func (f HandlerFunc) Handle(w io.Writer, req *Message) {
	f(w, req)
}

// Listener is an interface for receiving messages from a chat service like Messenger or Slack.
type Listener interface {
	Messages() <-chan *Message
	Close()
}

// Sender is an interface for sending messages to a chat service like Messenger or Slack.
//
// A Sender must be safe for concurrent use by multiple
// goroutines.
type Sender interface {
	Send(context.Context, *Response) error
	TypingOn(context.Context, User) error
	TypingOff(context.Context, User) error
	Close()
}

type matcher interface {
	MatchString(string) bool
}

type handlerEntry struct {
	matcher matcher
	handler Handler
}

// Bot receives events from a Listener, dispatches events to handlers, and sends
// responses back to a Sender.
type Bot struct {
	Listener      Listener
	Sender        Sender
	NumGoroutines int
	// NotFoundHandler will be called when no handlers match an incoming messaage.
	NotFoundHandler Handler

	mu       sync.Mutex
	handlers []*handlerEntry
	stopped  chan struct{}
	wg       sync.WaitGroup
}

func New() *Bot {
	return &Bot{
		NumGoroutines: 1,
		mu:            sync.Mutex{},
		handlers:      make([]*handlerEntry, 0),
		stopped:       make(chan struct{}),
		wg:            sync.WaitGroup{},
	}
}

func (b *Bot) Handle(pattern string, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	// Compile the pattern as a regular expression.
	re, err := regexp.Compile(fmt.Sprintf("(?is)%s", pattern))
	if err != nil {
		panic("botengine: invalid pattern " + pattern)
	}
	b.handlers = append(b.handlers, &handlerEntry{
		matcher: re,
		handler: handler,
	})
}

func (b *Bot) HandleFunc(pattern string, handler func(io.Writer, *Message)) {
	b.Handle(pattern, HandlerFunc(handler))
}

func (b *Bot) Run() {
	if b.Listener == nil {
		panic("botengine: Listener must not be nil")
	}
	if b.Sender == nil {
		panic("botengine: Sender must not be nil")
	}
	for i := 0; i < b.NumGoroutines; i++ {
		go b.run()
	}
}

func (b *Bot) run() {
	b.wg.Add(1)
	defer b.wg.Done()

	for {
		select {
		case msg := <-b.Listener.Messages():
			b.receive(msg)
		case <-b.stopped:
			return
		}
	}
}

func (b *Bot) receive(msg *Message) {
	for _, h := range b.handlers {
		if h.matcher.MatchString(msg.Body) {
			b.dispatch(h.handler, msg)
			return
		}
	}
	if b.NotFoundHandler != nil {
		b.dispatch(b.NotFoundHandler, msg)
	}
}

func (b *Bot) dispatch(handler Handler, msg *Message) {
	buf := &bytes.Buffer{}

	// Call handler.
	ctx := context.Background()
	_ = b.Sender.TypingOn(ctx, msg.Sender) // FIXME: Handle error.
	handler.Handle(buf, msg)
	_ = b.Sender.TypingOff(ctx, msg.Sender) // FIXME: Handle error.

	if body := buf.String(); body != "" {
		b.send(msg.Sender, body)
	}
}

func (b *Bot) send(recipient User, body string) {
	res := &Response{
		Recipient: recipient,
		Body:      body,
	}
	ctx := context.Background()
	_ = b.Sender.Send(ctx, res) // FIXME: Handle error.
}

func (b *Bot) Stop() {
	close(b.stopped)
	b.wg.Wait()

	b.Listener.Close()
	b.Sender.Close()
}
