package botengine

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"
)

const defaultNotFoundReply = "I don't understand 🤷"

type Source interface {
	Events() <-chan *Event
	Close()
}

type Sink interface {
	Flush(*Event) error
	Close()
}

type handler struct {
	// pattern is not case-sensitive
	pattern    string
	handleFunc func(io.Writer, *Event)
}

type Settings struct {
	NumGoroutines int
	NotFoundReply string
	Echo          bool
}

var DefaultSettings = Settings{
	NumGoroutines: 1,
	NotFoundReply: defaultNotFoundReply,
	Echo:          false,
}

// Engine provides the brain of a bot by dispatching events to handlers.
//
// type Bot struct {
//     *botengine.Engine
// }
type Engine struct {
	source Source
	sink   Sink

	// settings
	notFoundReply string
	numGoroutines int
	echo          bool

	mu       sync.Mutex
	handlers []*handler
	stopped  chan struct{}
	wg       sync.WaitGroup
}

func New(source Source, sink Sink, settings Settings) *Engine {
	return &Engine{
		source:        source,
		sink:          sink,
		notFoundReply: settings.NotFoundReply,
		numGoroutines: settings.NumGoroutines,
		echo:          settings.Echo,
		mu:            sync.Mutex{},
		handlers:      make([]*handler, 0),
		stopped:       make(chan struct{}),
		wg:            sync.WaitGroup{},
	}
}

func (e *Engine) Handle(pattern string, handleFunc func(io.Writer, *Event)) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.handlers = append(e.handlers, &handler{
		pattern:    strings.ToLower(pattern),
		handleFunc: handleFunc,
	})
}

func (e *Engine) Run() {
	for i := 0; i < e.numGoroutines; i++ {
		go e.run()
	}
}

func (e *Engine) run() {
	e.wg.Add(1)
	defer e.wg.Done()

	for {
		select {
		case ev := <-e.source.Events():
			e.dispatch(ev)
		case <-e.stopped:
			return
		}
	}
}

func (e *Engine) dispatch(ev *Event) {
	switch ev.Type {
	case MessageEvent:
		msg := ev.Object.(*Message)
		for _, h := range e.handlers {
			buf := &bytes.Buffer{}
			if h.pattern == strings.ToLower(msg.Text) {
				h.handleFunc(buf, ev)
				if reply := buf.String(); reply != "" {
					e.flush(msg.User, reply)
				}
				return
			}
			if e.echo {
				fmt.Fprintf(buf, "You sent the message \"%s\".", msg.Text)
				e.flush(msg.User, buf.String())
				return
			}
		}
		e.replyNotFound(msg.User)
	default:
		// Ignore unsupported events.
	}
}

func (e *Engine) flush(usr User, text string) {
	res := &Event{
		Type: MessageEvent,
		Object: &Message{
			User: usr,
			Text: text,
		},
	}
	_ = e.sink.Flush(res)
}

func (e *Engine) replyNotFound(usr User) {
	e.flush(usr, e.notFoundReply)
}

func (e *Engine) Stop() {
	close(e.stopped)
	e.wg.Wait()

	e.source.Close()
	e.sink.Close()
}
