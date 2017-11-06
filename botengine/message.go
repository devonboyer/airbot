package botengine

// Message represents an incoming message.
type Message struct {
	// Sender is the user who sent the message.
	Sender User
	// Body is the body of the message.
	Body string
}
