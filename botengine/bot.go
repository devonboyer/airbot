package botengine

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/pkg/errors"
)

// Message represents an incoming message.
type Message struct {
	// Sender is the user who sent the message.
	Sender User
	// Body is the body of the message.
	Body string
}

type Status int

const (
	// StatusOk indicates a bot is responsing to a message.
	StatusOk Status = iota
	// StatusError indicates a bot is responding with an error message.
	StatusError
	// StatusNotFound indicates a bot doesn't know how to handle a message.
	StatusNotFound
)

type Response struct {
	// Recipient is the user who should receive the message.
	Recipient User
	// Body is the body of the message.
	Body string
	// Status is the status of the response.
	Status Status
}

type ResponseWriter interface {
	io.Writer
	SetStatus(Status)
}

// User is a user in the chat service.
type User struct {
	// ID is the user's unique ID.
	ID string
}

type Handler interface {
	Handle(ResponseWriter, *Message)
}

type HandlerFunc func(ResponseWriter, *Message)

func (f HandlerFunc) Handle(w ResponseWriter, msg *Message) {
	f(w, msg)
}

func Error(w ResponseWriter, err error) {
	fmt.Fprintf(w, err.Error())
	w.SetStatus(StatusError)
}

func NotFound(w ResponseWriter, msg *Message) {
	fmt.Fprintf(w, fmt.Sprintf("I don't understand \"%s\".", msg.Body))
	w.SetStatus(StatusNotFound)
}

func NotFoundHandler() Handler { return HandlerFunc(NotFound) }

// ChatService is an interface for sending and receiving messages from a chat service
// like Messenger.
//
// A ChatService must be safe for concurrent use by multiple
// goroutines.
type ChatService interface {
	Messages() <-chan *Message
	Send(context.Context, *Response) error
	TypingOn(context.Context, User) error
	TypingOff(context.Context, User) error
	Close()
}

// Bot receives events from a ChatService, dispatches events to handlers via a Router, and sends
// responses back to the ChatService.
type Bot struct {
	ChatService ChatService
	Router      Router
	// NumGoroutines is the number of workers that can respond to incoming messages.
	// Default is 1.
	NumGoroutines int
	// NotFoundHandler will be called when no handlers match an incoming message.
	NotFoundHandler Handler

	stopped chan struct{}
	wg      sync.WaitGroup
}

func New() *Bot {
	return &Bot{
		Router:          DefaultRouter,
		NumGoroutines:   1,
		NotFoundHandler: NotFoundHandler(),
		stopped:         make(chan struct{}),
		wg:              sync.WaitGroup{},
	}
}

func (b *Bot) Run() <-chan error {
	if b.ChatService == nil {
		panic("botengine: ChatService must not be nil")
	}

	outError := make(chan error, 1)
	for i := 0; i < b.NumGoroutines; i++ {
		b.wg.Add(1)
		go b.run(outError)
	}
	return outError
}

func (b *Bot) run(outError chan error) {
	defer b.wg.Done()

	for {
		select {
		case msg := <-b.ChatService.Messages():
			b.receive(outError, msg)
		case <-b.stopped:
			return
		}
	}
}

func (b *Bot) receive(outError chan error, msg *Message) {
	handler := b.Router.Match(msg)
	if handler != nil {
		b.dispatch(outError, handler, msg)
	} else if b.NotFoundHandler != nil {
		b.dispatch(outError, b.NotFoundHandler, msg)
	}
}

func (b *Bot) dispatch(outError chan error, handler Handler, msg *Message) {
	rr := NewRecorder()

	// Call handler.
	ctx := context.Background()
	if err := b.ChatService.TypingOn(ctx, msg.Sender); err != nil {
		outError <- errors.Wrap(err, "could not send typing on action")
	}
	handler.Handle(rr, msg)
	if err := b.ChatService.TypingOff(ctx, msg.Sender); err != nil {
		outError <- errors.Wrap(err, "could not send typing off action")
	}
	if body := rr.Body.String(); body != "" {
		b.send(outError, msg.Sender, body, rr.Status)
	}
}

func (b *Bot) send(outError chan error, recipient User, body string, status Status) {
	res := &Response{
		Recipient: recipient,
		Body:      body,
		Status:    status,
	}
	ctx := context.Background()
	if err := b.ChatService.Send(ctx, res); err != nil {
		outError <- errors.Wrap(err, "could not send response")
	}
}

func (b *Bot) Stop() {
	close(b.stopped)
	b.wg.Wait()

	b.ChatService.Close()
}
