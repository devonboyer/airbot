package botengine

type Response struct {
	// Recipient is the user who should receive the message.
	Recipient User
	// Body is the body of the message.
	Body string
}
