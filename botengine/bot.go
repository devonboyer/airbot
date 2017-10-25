package botengine

import (
	"sync"
)

const (
	defaultNotFoundReply = "I don't understand 🤷"
	defaultErrorReply    = "Something went wrong"
)

type Option interface {
	Apply(*Bot)
}

func WithNotFoundReply(s string) Option {
	return withNotFoundReply(s)
}

type withNotFoundReply string

func (w withNotFoundReply) Apply(b *Bot) {
	b.notFoundReply = string(w)
}

func WithErrorReply(s string) Option {
	return withErrorReply(s)
}

type withErrorReply string

func (w withErrorReply) Apply(b *Bot) {
	b.errorReply = string(w)
}

type Event interface{} // TODO

type Source interface {
	Events() <-chan Event
	Close()
}

type Sink interface {
	Flush(Event) error
	Close()
}

type Bot struct {
	source Source
	sink   Sink

	notFoundReply, errorReply string

	mu       sync.Mutex
	handlers []*handler
	stopped  chan struct{}
	wg       sync.WaitGroup
}

type handler struct {
	pattern    string
	handleFunc func(string) (string, error)
}

func New(source Source, sink Sink, opts ...Option) *Bot {
	o := []Option{
		WithNotFoundReply(defaultNotFoundReply),
		WithErrorReply(defaultErrorReply),
	}
	opts = append(o, opts...)
	bot := &Bot{
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

func (b *Bot) Handle(pattern string, handleFunc func(string) (string, error)) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers = append(b.handlers, &handler{
		pattern:    pattern,
		handleFunc: handleFunc,
	})
}

func (b *Bot) Run() {
	go b.run()
}

func (b *Bot) run() {
	b.wg.Add(1)
	defer b.wg.Done()
}

func (b *Bot) Stop() {
	close(b.stopped)
	b.wg.Wait()

	b.source.Close()
	b.sink.Close()
}
