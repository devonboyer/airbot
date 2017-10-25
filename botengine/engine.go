package botengine

import (
	"sync"
)

const (
	defaultNotFoundReply = "I don't understand 🤷"
	defaultErrorReply    = "Something went wrong"
)

type Option interface {
	Apply(*Engine)
}

func WithNotFoundReply(s string) Option {
	return withNotFoundReply(s)
}

type withNotFoundReply string

func (w withNotFoundReply) Apply(b *Engine) {
	b.notFoundReply = string(w)
}

func WithErrorReply(s string) Option {
	return withErrorReply(s)
}

type withErrorReply string

func (w withErrorReply) Apply(b *Engine) {
	b.errorReply = string(w)
}

type Event interface{} // TODO

type Source interface {
	Events() <-chan *Event
	Close()
}

type Sink interface {
	Flush(*Event) error
	Close()
}

type handler struct {
	pattern    string
	handleFunc func(string) (string, error)
}

// Engine provides the brain of a bot by dispatching events to handlers.
//
// type Bot struct {
//     *botengine.Engine
// }
type Engine struct {
	source Source
	sink   Sink

	notFoundReply, errorReply string

	mu       sync.Mutex
	handlers []*handler
	stopped  chan struct{}
	wg       sync.WaitGroup
}

func New(source Source, sink Sink, opts ...Option) *Engine {
	o := []Option{
		WithNotFoundReply(defaultNotFoundReply),
		WithErrorReply(defaultErrorReply),
	}
	opts = append(o, opts...)
	bot := &Engine{
		source:   source,
		sink:     sink,
		mu:       sync.Mutex{},
		handlers: make([]*handler, 0),
		stopped:  make(chan struct{}),
		wg:       sync.WaitGroup{},
	}
	for _, opt := range opts {
		opt.Apply(bot)
	}
	return bot
}

func (b *Engine) Handle(pattern string, handleFunc func(string) (string, error)) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers = append(b.handlers, &handler{
		pattern:    pattern,
		handleFunc: handleFunc,
	})
}

func (b *Engine) Run() {
	go b.run()
}

func (b *Engine) run() {
	b.wg.Add(1)
	defer b.wg.Done()

	for {
		select {
		case ev := <-b.source.Events():
			b.dispatch(ev)
		case <-b.stopped:
			return
		}
	}
}

func (b *Engine) dispatch(ev *Event) {
	// Handle event
}

func (b *Engine) Stop() {
	close(b.stopped)
	b.wg.Wait()

	b.source.Close()
	b.sink.Close()
}
