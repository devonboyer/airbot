package botengine

type EventType string

const MessageEvent EventType = "message"

// Input
type Event struct {
	Type   EventType
	Object interface{}
}

type User struct {
	ID string
}

// Output
type Message struct {
	User User
	Text string
}

type Response struct {
	// The pattern the incoming message matched.
	Pattern string
	// An array of arguments created by shellsplitting the message body, as if
	// it were a shell command.
	Args []string
}

// Handler.respond_to?
// handler.dispatch
