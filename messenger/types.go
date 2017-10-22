package messenger

type Recipient struct {
	ID string `json:"id"`
}

type Message struct {
	Text string `json:"text,omitempty"`
}

type SendMessageBody struct {
	Recipient Recipient `json:"recipient"`
	Message   Message   `json:"message"`
	NotifType string    `json:"notification_type"`
}

type SenderActionBody struct {
	Recipient Recipient `json:"recipient"`
	Action    string    `json:"sender_action"`
}

type WebhookEvent struct {
	Object  string `json:"object"`
	Entries []struct {
		PageID    int64         `json:"id"`
		Time      int64         `json:"time"`
		Messaging []interface{} `json:"messaging"`
	} `json:"entry"`
}

type MessageEvent struct {
	Sender struct {
		ID string `json:"id"`
	} `json:"sender"`
	Recipient struct {
		ID string `json:"id"`
	} `json:"recipient"`
	Timestamp int `json:"timestamp,omitempty"`
	Message   struct {
		MID  string `json:"mid"`
		Text string `json:"text,omitempty"`
	} `json:"message"`
}
