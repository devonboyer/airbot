package messenger

type Recipient struct {
	ID string `json:"id"`
}

type Message struct {
	Text string `json:"text,omitempty"`
}

type MessageRequest struct {
	Recipient Recipient `json:"recipient"`
	Message   Message   `json:"message"`
	NotifType string    `json:"notification_type"`
}
