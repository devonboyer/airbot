package botengine

type EventType string

const MessageEvent EventType = "message"

type Event struct {
	Type   EventType
	Object interface{}
}

type User struct {
	ID string
}

type Message struct {
	User User
	Text string
}
