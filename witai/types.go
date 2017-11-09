package witai

type Message struct {
	MsgID    string                   `json:"msg_id"`
	Text     string                   `json:"_text"`
	Entities map[string][]interface{} `json:"entities"`
}
