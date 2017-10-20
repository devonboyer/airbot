package raw

type Request struct {
	Recipient struct {
		ID string `json:"id,omitempty"`
	} `json:"recipient"`
	Message struct {
		Text string `json:"text,omitempty"`
	} `json:"message,omitempty"`
	Action    string `json:"sender_action,omitempty"`
	NotifType string `json:"notification_type"`
}
