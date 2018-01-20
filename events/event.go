package events

type Receiver interface {
  Receive(payload interface {})
}

type Event struct {
  timestamp int
  payload   interface {}
  receiver  *Receiver
}

func NewEvent(timestamp int, payload interface {}, receiver *Receiver) *Event {
  e := new(Event)

  e.timestamp = timestamp
  e.payload   = payload
  e.receiver  = receiver

  return e
}

func (e *Event) Timestamp() int {
  return e.timestamp
}
