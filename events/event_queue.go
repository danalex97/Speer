package events

import (
  . "github.com/danalex97/Speer/structs"
  "sync"
  "fmt"
)

type EventQueue interface {
  Push(e *Event)
  Pop() *Event
}

const LazyQueueChanSize int = 100

type lazyEventQueue struct {
  pq PriorityQueue
  stream chan *Event

  *sync.Mutex
}

func NewLazyEventQueue() EventQueue {
  eq := new(lazyEventQueue)

  eq.pq     = NewPriorityQueue()
  eq.stream = make(chan *Event, LazyQueueChanSize + 5)
  eq.Mutex  = new(sync.Mutex)

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
    // We need to avoid multiple depressure operations at the same time.
    eq.Lock()
    defer eq.Unlock()

    // it must be that the channel is full, so we need to
    // relase some pressure
    eq.depressure()
    eq.pq.Push(Int(event.timestamp), event)
  }
}

func (eq *lazyEventQueue) Pop() (event *Event) {
  eq.Lock()
  defer eq.Unlock()

  done := false
  for !done {
    select {
    case ev := <-eq.stream:
      eq.pq.Push(Int(ev.timestamp), ev)
    default:
      done = true
    }
  }

  if eq.pq.Len() > 0 {
    event = eq.pq.Pop().Value.(*Event)
  }
  return
}
