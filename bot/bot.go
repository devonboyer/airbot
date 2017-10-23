package bot

import (
	"fmt"
	"sync"
)

const (
	defaultNotFoundReply = "I don't understand 🤷"
	defaultErrorReply    = "Something went wrong"
)

type Message struct {
	SenderID string
	Text     string
}

type Reply struct {
	RecipientID string
	Text        string
}

type Listener interface {
	Messages() <-chan Message
}

type Sender interface {
	Send(Reply)
}

type SenderFunc func(Reply)

func (f SenderFunc) Send(r Reply) {
	f(r)
}

type Bot struct {
	NotFoundReply string
	ErrorReply    string

	listener Listener
	sender   Sender
	mu       sync.Mutex
	handlers []*handler
	stopped  chan struct{}
	wg       sync.WaitGroup
}

type handler struct {
	pattern    string
	handleFunc func(string) (string, error)
}

func New(listener Listener, sender Sender) *Bot {
	return &Bot{
		NotFoundReply: defaultNotFoundReply,
		ErrorReply:    defaultErrorReply,
		listener:      listener,
		sender:        sender,
		mu:            sync.Mutex{},
		handlers:      make([]*handler, 0),
		stopped:       make(chan struct{}),
		wg:            sync.WaitGroup{},
	}
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
	for {
		select {
		case msg := <-b.listener.Messages():
			b.dispatch(msg)
		case <-b.stopped:
			return
		}
	}
}

func (b *Bot) dispatch(msg Message) {
	id := msg.SenderID
	for _, h := range b.handlers {
		// For now, use very simplistic string comparison to dispatch to correct handler.
		if h.pattern == msg.Text {
			msg, err := h.handleFunc(msg.SenderID)
			if err != nil {
				b.reply(id, b.ErrorReply)
				fmt.Println(err)
				return
			}
			b.reply(id, msg)
			return
		}
	}
	b.reply(id, b.NotFoundReply)
}

func (b *Bot) reply(recipientID, msg string) {
	reply := Reply{
		RecipientID: recipientID,
		Text:        msg,
	}
	b.sender.Send(reply)
}

func (b *Bot) Stop() {
	close(b.stopped)
	b.wg.Wait()
}
