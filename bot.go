package airbot

type Command struct {
	Pattern     string
	Help        string
	Subcommands []Command
	Handler     func() error
}

type Receiver interface {
	Receive() <-chan string
}

type Sender interface {
	Receive() <-chan string
	Send(string)
	TypingOn()
	TypingOff()
}

type Bot struct {
	Receiver Receiver
	Sender   Sender
	commands []Command
}

func NewBot(receiver Receiver, sender Sender) *Bot {
	return &Bot{
		Receiver: receiver,
		Sender:   sender,
		commands: make([]Command, 0),
	}
}

func (b *Bot) AddCommand(cmd Command) {
	b.commands = append(b.commands, cmd)
}

func (b *Bot) Run() <-chan error {

}
