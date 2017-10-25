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

func WithNumGoroutines(n int) Option {
	return withNumGoroutines(n)
}

type withNumGoroutines int

func (w withNumGoroutines) Apply(e *Engine) {
	e.numGoroutines = int(w)
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

	// options
	notFoundReply string
	errorReply    string
	numGoroutines int

	mu       sync.Mutex
	handlers []*handler
	stopped  chan struct{}
	wg       sync.WaitGroup
}

func New(source Source, sink Sink, opts ...Option) *Engine {
	o := []Option{
		WithNumGoroutines(1),
		WithNotFoundReply(defaultNotFoundReply),
		WithErrorReply(defaultErrorReply),
	}
	opts = append(o, opts...)
	eng := &Engine{
		source:   source,
		sink:     sink,
		mu:       sync.Mutex{},
		handlers: make([]*handler, 0),
		stopped:  make(chan struct{}),
		wg:       sync.WaitGroup{},
	}
	for _, opt := range opts {
		opt.Apply(eng)
	}
	return eng
}

func (e *Engine) Handle(pattern string, handleFunc func(string) (string, error)) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.handlers = append(e.handlers, &handler{
		pattern:    pattern,
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
	// Handle event
}

func (e *Engine) Stop() {
	close(e.stopped)
	e.wg.Wait()

	e.source.Close()
	e.sink.Close()
}
