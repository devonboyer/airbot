package botengine

import "context"

type Event interface{}

type Queue interface {
	Push(context.Context, Event)
	Pop(context.Context) Event
	Close()
}

const queueBufferSize = 1024

// Very simplistic queue implementation. In reality, if we try to buffer events in memory
// process will be OOMkilled or be forced to drop events.
type queue chan Event

func NewQueue() Queue {
	return make(queue, queueBufferSize)
}

func (q queue) Push(ctx context.Context, ev Event) {
	select {
	case q <- ev:
	default: // Ignore events if the buffer is full :(
	}
}

func (q queue) Pop(ctx context.Context) Event {
	select {
	case ev := <-q:
		return ev
	}
}

func (q queue) Close() {}
