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
	Handle(io.Writer, *Event)
}

type HandlerFunc func(io.Writer, *Event)

func (f HandlerFunc) Handle(w io.Writer, ev *Event) {
	f(w, ev)
}

type matcher interface {
	MatchString(string) bool
}

type handlerEntry struct {
	matcher matcher
	handler Handler
}

type Settings struct {
	NumGoroutines int
	NotFoundReply string
	// ShouldEcho will echo responses to the Chatter if no handler matches the received event.
	// This may be useful for debugging. Defaults to false.
	ShouldEcho bool
}

var DefaultSettings = Settings{
	NumGoroutines: 1,
	NotFoundReply: "I don't understand 🤷",
	ShouldEcho:    false,
}

// Listener is an interface for receiving messages from a chat service like Messenger or Slack.
type Listener interface {
	Events() <-chan *Event
	Close()
}

// Sender is an interface for sending messages to a chat service like Messenger or Slack.
//
// A Sender must be safe for concurrent use by multiple
// goroutines.
type Sender interface {
	Send(context.Context, *Event) error
	Close()
}

// Bot receives events from a Listener, dispatches events to handlers, and sends
// responses back to a Sender.
type Bot struct {
	Listener Listener
	Sender   Sender

	// settings
	notFoundReply string
	numGoroutines int
	shouldEcho    bool

	mu       sync.Mutex
	handlers []*handlerEntry
	stopped  chan struct{}
	wg       sync.WaitGroup
}

func New(settings Settings) *Bot {
	return &Bot{
		notFoundReply: settings.NotFoundReply,
		numGoroutines: settings.NumGoroutines,
		shouldEcho:    settings.ShouldEcho,
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

func (b *Bot) HandleFunc(pattern string, handler func(io.Writer, *Event)) {
	b.Handle(pattern, HandlerFunc(handler))
}

func (b *Bot) Run() {
	if b.Listener == nil {
		panic("botengine: Listener must not be nil")
	}
	if b.Sender == nil {
		panic("botengine: Sender must not be nil")
	}
	for i := 0; i < b.numGoroutines; i++ {
		go b.run()
	}
}

func (b *Bot) run() {
	b.wg.Add(1)
	defer b.wg.Done()

	for {
		select {
		case ev := <-b.Listener.Events():
			b.dispatch(ev)
		case <-b.stopped:
			return
		}
	}
}

func (b *Bot) dispatch(ev *Event) {
	switch ev.Type {
	case MessageEvent:
		msg := ev.Object.(*Message)
		for _, h := range b.handlers {
			buf := &bytes.Buffer{}
			if h.matcher.MatchString(msg.Text) {
				h.handler.Handle(buf, ev)
				if reply := buf.String(); reply != "" {
					b.send(msg.User, reply)
				}
				return
			}
			if b.shouldEcho {
				fmt.Fprintf(buf, "You sent the message \"%s\".", msg.Text)
				b.send(msg.User, buf.String())
				return
			}
		}
		b.replyNotFound(msg.User)
	default:
		// Ignore unsupported events.
	}
}

func (b *Bot) send(usr User, text string) {
	res := &Event{
		Type: MessageEvent,
		Object: &Message{
			User: usr,
			Text: text,
		},
	}
	// FIXME: This error should be handled more gracefully.
	ctx := context.Background()
	_ = b.Sender.Send(ctx, res)
}

func (b *Bot) replyNotFound(usr User) {
	b.send(usr, b.notFoundReply)
}

func (b *Bot) Stop() {
	close(b.stopped)
	b.wg.Wait()

	b.Listener.Close()
	b.Sender.Close()
}
