package events

type Receiver interface {
	Receive(*Event) *Event
}

type Producer interface {
	Produce() *Event
}
