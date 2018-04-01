package events

import (
  "runtime"
)

type EventObserver interface {
  EventChan() <-chan *Event
  EnqueEvent(*Event)
}

const maxObserverQueue int = 1000

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

type globalObserver struct {
  observer chan *Event
}

func NewGlobalEventObserver() EventObserver {
  obs := new(globalObserver)
  obs.observer = make(chan *Event, maxObserverQueue)
  return obs
}

func (o *globalObserver) EventChan() <-chan *Event {
  return o.observer
}

func (o *globalObserver) EnqueEvent(e *Event) {
  select {
  case o.observer <- e:
  default:
    // If we can't register the observer, let someone else to run
    // to break the livelock.
    runtime.Gosched()
    o.EnqueEvent(e)
  }
}
