package events

import (
  "fmt"
  "runtime"
  "sync"
)

type Simulation struct {
  newObservers chan EventObserver
  observers []EventObserver
  stopped chan interface {}
  timeMutex *sync.RWMutex
  time    int
  EventQueue
}

func NewLazySimulation() (s Simulation) {
  s = Simulation{
    make(chan EventObserver, 50),
    make([]EventObserver, 0),
    make(chan interface {}, 1),
    new(sync.RWMutex),
    0,
    NewLazyEventQueue(),
  }
  return
}

func (s *Simulation) RegisterObserver(eventObserver EventObserver) {
  s.newObservers <- eventObserver
}

func (s *Simulation) Time() int {
  s.timeMutex.RLock()
  defer s.timeMutex.RUnlock()

  return s.time
}

func (s *Simulation) Stop() {
  s.stopped <- nil
}

func (s *Simulation) Run() {
  fmt.Println("Starting the simulation.")

  for {
    select {
    case <-s.stopped:
      break
    case observer := <-s.newObservers:
      // fmt.Println("New Observer >", observer)

      s.observers = append(s.observers, observer)
    default:
      if event:= s.Pop(); event != nil {
        // fmt.Println("Event received >", event)

        // The event gets dispached to observers
        for _, observer := range(s.observers) {
          observer.EnqueEvent(event)
        }

        s.timeMutex.Lock()
        s.time = event.timestamp
        s.timeMutex.Unlock()

        receiver := event.receiver

        if receiver == nil {
          continue
        }

        newEvent := receiver.Receive(event)

        if newEvent != nil {
          s.Push(newEvent)
        }
      } else {
        runtime.Gosched()
      }
    }
  }
}
