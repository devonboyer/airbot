package messenger

type Recipient struct {
	ID string `json:"id"`
}

type Message struct {
	MID  string `json:"mid"`
	Text string `json:"text,omitempty"`
}

type SendBody struct {
	Recipient Recipient `json:"recipient"`
	Message   Message   `json:"message"`
	NotifType string    `json:"notification_type"`
}

type SenderActionBody struct {
	Recipient Recipient `json:"recipient"`
	Action    string    `json:"sender_action"`
}

type ReceiveBody struct {
	Sender    Recipient `json:"sender"`
	Recipient Recipient `json:"recipient"`
	Timestamp int       `json:"timestamp"`
	Message   Message   `json:"message"`
}
