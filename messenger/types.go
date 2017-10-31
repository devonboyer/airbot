package messenger

type Recipient struct {
	ID string `json:"id"`
}

type Message struct {
	Text string `json:"text,omitempty"`
}

type MessageMarshaler struct {
	Recipient Recipient `json:"recipient"`
	Message   Message   `json:"message"`
	NotifType string    `json:"notification_type"`
}

type SenderActionMarshaler struct {
	Recipient Recipient `json:"recipient"`
	Action    string    `json:"sender_action"`
}

type Event struct {
	Object  string `json:"object"`
	Entries []struct {
		PageID    string     `json:"id"`
		Time      int64      `json:"time"`
		Messaging []Callback `json:"messaging"`
	} `json:"entry"`
}

type Callback struct {
	Sender struct {
		ID string `json:"id"`
	} `json:"sender"`
	Recipient struct {
		ID string `json:"id"`
	} `json:"recipient"`
	Timestamp int64 `json:"timestamp"`
	Message   *struct {
		MID  string `json:"mid"`
		Seq  int64  `json:"seq"`
		Text string `json:"text,omitempty"`
	} `json:"message,omitempty"`
}
