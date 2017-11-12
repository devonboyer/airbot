package botengine

import (
	"strings"
	"sync"
)

type Router interface {
	Match(*Message) Handler
}

type CommandRouter struct {
	mu   sync.Mutex
	cmds []*command
}

func NewCommandRouter() *CommandRouter {
	return &CommandRouter{
		mu:   sync.Mutex{},
		cmds: make([]*command, 0),
	}
}

type command struct {
	pattern string
	handler Handler
}

func (r *CommandRouter) Match(msg *Message) Handler {
	for _, h := range r.cmds {
		if strings.ToLower(h.pattern) == msg.Body {
			return h.handler
		}
	}
	return nil
}

func (r *CommandRouter) Handle(pattern string, handler Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cmds = append(r.cmds, &command{
		pattern: pattern,
		handler: handler,
	})
}

func (r *CommandRouter) HandleFunc(pattern string, handler func(ResponseWriter, *Message)) {
	r.Handle(pattern, HandlerFunc(handler))
}

var DefaultRouter = NewCommandRouter()

func Handle(pattern string, handler Handler) {
	DefaultRouter.Handle(pattern, handler)
}

func HandleFunc(pattern string, handler func(ResponseWriter, *Message)) {
	DefaultRouter.HandleFunc(pattern, handler)
}
