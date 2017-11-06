package botengine

import "strings"

// Message represents an incoming message.
type Message struct {
	// Sender is the user who sent the message.
	Sender User
	// Body is the body of the message.
	Body string
}

// Args returns an array of arguments created by shellsplitting the message body, as if
// it were a shell command.
func (m *Message) Args() []string {
	return strings.Split(m.Body, " ")
}
