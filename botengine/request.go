package botengine

type Request struct {
	// Message is the incoming message
	Message *Message
	// Args is an array of arguments created by shellsplitting the message body, as if
	// it were a shell commands
	Args []string
}
