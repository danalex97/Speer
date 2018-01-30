package events

import (
  "fmt"
  "sync"
)

type Simulation struct {
  sync.Mutex

  observers []EventObserver
  stopped chan interface {}
  time    int
  EventQueue
}

func NewLazySimulation() (s Simulation) {
  s = Simulation{
    *new(sync.Mutex),
    make([]EventObserver, 0),
    make(chan interface {}),
    0,
    NewLazyEventQueue(),
  }
  return
}

func (s *Simulation) RegisterObserver(eventObserver EventObserver) {
  s.Lock()
  defer s.Unlock()

  s.observers = append(s.observers, eventObserver)
}

func (s *Simulation) Time() int {
  s.Lock()
  defer s.Unlock()

  return s.time
}

func (s *Simulation) Stop() {
  s.Lock()
  defer s.Unlock()

  s.stopped <- nil
}

func (s *Simulation) Run() {
  s.time = 0
  for {
    select {
    case <-s.stopped:
      break
    default:
      s.Lock()
      defer s.Unlock()

      if event:= s.Pop(); event != nil {
        fmt.Println("Event received >", event.timestamp)

        // The event gets dispached to observers
        for _, observer := range(s.observers) {
          observer.EnqueEvent(event)
        }

        s.time = event.timestamp
        receiver := event.receiver

        if receiver == nil {
          continue
        }

        newEvent := receiver.Receive(event)

        if newEvent != nil {
          s.Push(newEvent)
        }
      }
    }
  }
}
