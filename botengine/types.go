package botengine

const MessageEvent = "message"

type Event struct {
	Kind   string
	Object interface{}
}

type User struct {
	ID string
}

type Message struct {
	User User
	Text string
}
