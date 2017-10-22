package airbot

type Event interface{}

type Listener interface {
	Events() <-chan Event
}

type Translator func(string) string

type Handler struct {
	pattern    string
	translator Translator
}

type Responder interface {
	Respond(string) error
}

type Bot struct {
	Listener Listener
	handlers []*Handler
}

func NewBot() *Bot {
	return &Bot{
		handlers: make([]*Handler, 0),
	}
}

func (b *Bot) Handle(pattern string, translator Translator) {
	handler := &Handler{
		pattern:    pattern,
		translator: translator,
	}
	b.handlers = append(b.handlers, handler)
}

func (b *Bot) Run() {
	for {
		select {
		case event := <-b.Listener.Events():
			b.dispatch(event)
		}
	}
}

func (b *Bot) dispatch(event Event) {
	// Check if I have a hook that matches this event
}
