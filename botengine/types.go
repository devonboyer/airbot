package botengine

type EventType string

const 

type Object interface {
	GetObjectKind() string
}

type Event struct {
	Entries []Object
}

type Message interface {
	// thing you give me

	GetMessage() string
	//
	// Sender/Recipient
	// UID
	//
	// sender ID, message body
}

func (msg *Message) GetObjectKind() string {

}

// ask this object what kind it is?
// botengine is the one doing the casting
