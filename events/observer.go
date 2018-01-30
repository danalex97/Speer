package events

type EventObserver interface {
  EventChan() <-chan *Event
  EnqueEvent(*Event)
}

type eventObserver struct {
  receiver Receiver
  observer chan *Event
}

func NewEventObserver(receiver Receiver) EventObserver {
  obs := new(eventObserver)
  obs.observer = make(chan *Event)
  obs.receiver = receiver
  return obs
}

func (o *eventObserver) EventChan() <-chan *Event {
  return o.observer
}

func (o *eventObserver) EnqueEvent(e *Event) {
  if e.Receiver() == o.receiver {
    o.observer <- e
  }
}

func NewGlobalEventObserver(receiver Receiver) EventObserver {
  obs := new(eventObserver)
  obs.observer = make(chan *Event)
  return obs
}

type globalObserver struct {
  observer chan *Event
}

func (o *globalObserver) EventChan() <-chan *Event {
  return o.observer
}

func (o *globalObserver) EnqueEvent(e *Event) {
  o.observer <- e
}
