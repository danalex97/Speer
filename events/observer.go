package events

import (
  "runtime"
)

type Observer interface {
  Recv() <-chan interface{}
  EnqueEvent(*Event)
}

const maxGlobalObserverQueue int = 1000
const maxObserverQueue       int = 1000

type eventObserver struct {
  receiver Receiver
  observer chan interface {}
}

func NewEventObserver(receiver Receiver) Observer {
  return &eventObserver{
    observer : make(chan interface {}, maxObserverQueue),
    receiver : receiver,
  }
}

func (o *eventObserver) Recv() <-chan interface{} {
  return o.observer
}

func (o *eventObserver) EnqueEvent(e *Event) {
  if e.Receiver() == o.receiver {
    o.observer <- e
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
