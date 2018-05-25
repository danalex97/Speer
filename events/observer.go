package events

import (
  "runtime"
)

type Observer interface {
  Recv() <-chan interface{}
  EnqueEvent(*Event)
}

type DecorableObserver interface {
  Decorable
  Observer
}

const maxGlobalObserverQueue int = 1000
const maxObserverQueue       int = 1000

type EventObserver struct {
  *Decorator

  receiver Receiver
  observer chan interface {}
}

func NewEventObserver(receiver Receiver) *EventObserver {
  return &EventObserver{
    Decorator : NewDecorator(),

    observer  : make(chan interface {}, maxObserverQueue),
    receiver  : receiver,
  }
}

func (o *EventObserver) Recv() <-chan interface{} {
  return o.observer
}

func (o *EventObserver) EnqueEvent(e *Event) {
  if e.Receiver() == o.receiver {
    o.observer <- o.Proxy(e)
  }
}

type globalObserver struct {
  observer chan interface {}
}

func NewGlobalEventObserver() Observer {
  return &globalObserver{
    observer : make(chan interface {}, maxObserverQueue),
  }
}

func (o *globalObserver) Recv() <-chan interface {} {
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
