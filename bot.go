package airbot

type Event struct {
	SenderID string
	Message  string
}

type Reply struct {
	RecipientID string
	Message     string
}

type Listener interface {
	Events() <-chan Event
}

type Handler func(string) (string, error)

type Command struct {
	pattern string
	handler Handler
}

type Responder interface {
	Respond(Reply)
}

type Bot struct {
	Listener  Listener
	Responder Responder
	commands  []*Command
}

func NewBot() *Bot {
	return &Bot{
		commands: make([]*Command, 0),
	}
}

func (b *Bot) Handle(pattern string, handler Handler) {
	cmd := &Command{
		pattern: pattern,
		handler: handler,
	}
	b.commands = append(b.commands, cmd)
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
	recipientID := event.SenderID
	for _, cmd := range b.commands {
		// For now, use very simplistic string comparison to dispatch to correct handler.
		if cmd.pattern == event.Message {
			msg, err := cmd.handler(event.SenderID)
			if err != nil {
				b.reply(recipientID, "Something went wrong")
				return
			}
			b.reply(recipientID, msg)
			return
		}
	}
	b.reply(recipientID, "No command found")
}

func (b *Bot) reply(recipientID, msg string) {
	reply := Reply{
		RecipientID: recipientID,
		Message:     msg,
	}
	b.Responder.Respond(reply)
}
