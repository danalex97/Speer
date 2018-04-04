package events

import (
  . "github.com/danalex97/Speer/structs"
  "fmt"
)

type EventQueue interface {
  Push(e *Event)
  Pop() *Event
}

const LazyQueueChanSize int = 50

type lazyEventQueue struct {
  pq PriorityQueue
  stream chan *Event
}

func NewLazyEventQueue() EventQueue {
  eq := new(lazyEventQueue)
  eq.pq     = NewPriorityQueue()
  eq.stream = make(chan *Event, LazyQueueChanSize + 5)
  return eq
}

func (eq *lazyEventQueue) depressure() {
  fmt.Println("Priority queue push channel full. Clearing it.")
  for len(eq.stream) > LazyQueueChanSize / 3 {
    event := <- eq.stream
    eq.pq.Push(Int(event.timestamp), event)
  }
}

func (eq *lazyEventQueue) Push(event *Event) {
  select {
  case eq.stream <- event:
  default:
    // it must be that the channel is full, so we need to
    // relase some pressure
    eq.depressure()
    eq.pq.Push(Int(event.timestamp), event)
  }
}

func (eq *lazyEventQueue) Pop() (event *Event) {
  eq.Push(nil)
  for {
    event = <- eq.stream
    if event == nil {
      break
    }

    eq.pq.Push(Int(event.timestamp), event)
  }

  if eq.pq.Len() > 0 {
    event = eq.pq.Pop().Value.(*Event)
  }
  return
}
