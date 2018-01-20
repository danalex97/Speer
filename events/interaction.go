package events

type Receiver interface {
  Receive(*Event)
}

type Producer interface {
  Produce() *Event
}
