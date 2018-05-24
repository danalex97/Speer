package events

import (
  "reflect"
  "runtime"
  "sync"
  "fmt"
)

type Simulation struct {
  newObservers chan EventObserver
  observers []EventObserver
  stopped chan interface {}
  timeMutex *sync.RWMutex
  time    int

  parallel bool

  EventQueue
}

const maxRegisterQueue int = 50
const minRegisterQueue int = 10

func NewLazySimulation() (s Simulation) {
  s = Simulation{
    newObservers : make(chan EventObserver, maxRegisterQueue),
    observers    : make([]EventObserver, 0),
    stopped      : make(chan interface {}, 1),
    timeMutex    : new(sync.RWMutex),
    time         : 0,

    parallel     : true,
    EventQueue   : NewLazyEventQueue(),
  }
  return
}

func (s *Simulation) SetParallel(parallel bool) {
  s.parallel = parallel
}

func (s *Simulation) RegisterProgress(property *ProgressProperty) {
  event := NewEvent(s.Time(), nil, property)
  s.Push(event)
}

func (s *Simulation) RegisterObserver(eventObserver EventObserver) {
  select {
  case s.newObservers <- eventObserver:
  default:
    // The observer queue is full, so we need to register the new observers
    // to clean in.
    for len(s.newObservers) > minRegisterQueue {
      observer := <-s.newObservers
      s.observers = append(s.observers, observer)
    }
    s.RegisterObserver(eventObserver)
  }
}

func (s *Simulation) Time() int {
  s.timeMutex.RLock()
  defer s.timeMutex.RUnlock()

  return s.time
}

func (s *Simulation) Stop() {
  s.stopped <- nil
}

func (s *Simulation) processEvent(event *Event) {
  newEvent := event.receiver.Receive(event)

  if newEvent != nil {
    s.Push(newEvent)
  }
}

func (s *Simulation) Run() {
  fmt.Println("Starting the simulation.")

  handler := s.Handle
  if s.parallel {
    handler = s.HandleParallel
  }

  for {
    select {
    case <-s.stopped:
      break
    case observer := <-s.newObservers:
      // fmt.Println("New Observer >", observer)

      s.observers = append(s.observers, observer)
    default:
      if s.Time() < 200 {
        // At the beginning we run the sequential simulator.
        s.Handle()
      } else {
        handler()
      }
    }
  }
}

// Handling events synchronously
func (s *Simulation) Handle() {
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
      return
    }

    newEvent := receiver.Receive(event)

    if newEvent != nil {
      s.Push(newEvent)
    }
  } else {
    runtime.Gosched()
  }

  return
}

func (s *Simulation) processGroup(group []*Event, done chan bool) {
  toProcess := []*Event{}
  for _, event := range group {
    toProcess = append(toProcess, event)
  }

  go func() {
    for _, event := range toProcess {
      s.processEvent(event)
    }
    done <- true
  }()
}

// Handling events in parallel
func (s *Simulation) HandleParallel() {
  event := s.Pop()
  if event == nil {
    // Empty event queue, so we can let other threads run.
    runtime.Gosched()
    return
  }

  // Starting new timeslice parallel execution.
  eventTime := event.timestamp
  // fmt.Println("New slice", eventTime)

  // Get all events happening at the same time for processing.
  events := []*Event{}
  for ;event != nil && event.timestamp == eventTime; event = s.Pop() {
    events = append(events, event)
  }
  // Push back the next event extracted.
  if event != nil {
    s.Push(event)
  }

  // The event gets dispached to observers.
  for _, event := range(events) {
    for _, observer := range(s.observers) {
      observer.EnqueEvent(event)
    }
  }

  // Modify the global timestamp
  s.timeMutex.Lock()
  s.time = events[0].timestamp
  s.timeMutex.Unlock()

  // Find special events.
  special := []*Event{}
  rest    := []*Event{}
  for _, event := range events {
    if event == nil || event.Receiver() == nil {
      continue
    }
    if reflect.TypeOf(event.Receiver()).String() == "*underlay.shortestPathRouter" {
      rest = append(rest, event)
    } else {
      special = append(special, event)
    }
  }
  events = rest

  // Process special events
  for _, event := range special {
    s.processEvent(event)
  }

  // Group events by Receiver.
  groups := make(map[Receiver][]*Event)
  for _, event := range(events) {
    receiver := event.Receiver()

    if receiver == nil {
      continue
    }

    if _, ok := groups[receiver]; !ok {
      groups[receiver] = []*Event{}
    }
    groups[receiver] = append(groups[receiver], event)
  }

  // Exectue events in parallel for each receiver group and wait for
  // responses.
  switch {
  case len(groups) == 0:
  case len(groups) == 1:
    // If there is only one group run events sequentially.
    for _, group := range groups {
      for _, event := range group {
        s.processEvent(event)
      }
    }
  default:
    fmt.Println("Parallel")
    for _, group := range groups {
      fmt.Printf("Group > %p \n", group[0].Receiver())
      fmt.Println("  Receiver type > ", reflect.TypeOf(group[0].Receiver()))
      for _, event := range group {
        fmt.Println("   Event received >", event, reflect.TypeOf(event.Payload()))
      }
    }

    done     := make(chan bool)
    routines := 0
    for _, group := range groups {
      // Process group.
      routines += 1
      s.processGroup(group, done)
    }

    for i := 0; i < routines; i++ {
      <-done
    }
  }

  return
}
