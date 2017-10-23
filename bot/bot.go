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

type Source interface {
	Messages() <-chan Message
	Send(Reply)
}

type Bot struct {
	NotFoundReply string
	ErrorReply    string

	source   Source
	mu       sync.Mutex
	handlers []*handler
	stopped  chan struct{}
	wg       sync.WaitGroup
}

type handler struct {
	pattern    string
	handleFunc func(string) (string, error)
}

func New(source Source) *Bot {
	return &Bot{
		NotFoundReply: defaultNotFoundReply,
		ErrorReply:    defaultErrorReply,
		source:        source,
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
		case msg := <-b.source.Messages():
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
	b.source.Send(reply)
}

func (b *Bot) Stop() {
	close(b.stopped)
	b.wg.Wait()
}
