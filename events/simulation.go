package events

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"
)

// The Simulation works as follows. At each moment the event with the
// smallest timestamp is popped out of the priority queue. When we pop the
// element we also notify all the observers associated with the event’s
// receiver. Registering an observer has priority over processing events, thus
// allowing the user to register observers at any moment. Since we allow any
// process to read the time, updating the time with the time of the current
// event is protected via a read-write lock.
type Simulation struct {
	observers    []Observer
	stopped      chan interface{}
	timeMutex    *sync.RWMutex
	time         int

	parallel bool
	check    func(Receiver) bool

	EventQueue
}

const maxRegisterQueue int = 50
const minRegisterQueue int = 10
const settleTime int = 200

func check(r Receiver) bool {
	return reflect.TypeOf(r).String() == "*underlay.shortestPathRouter"
}

func NewLazySimulation() (s *Simulation) {
	s = &Simulation{
		observers:    make([]Observer, 0),
		stopped:      make(chan interface{}, 1),
		timeMutex:    new(sync.RWMutex),
		time:         0,
		check:        check,

		parallel:   false,
		EventQueue: NewLazyEventQueue(),
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

func (s *Simulation) RegisterObserver(eventObserver Observer) {
	s.observers = append(s.observers, eventObserver)
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

	settled := false

	for {
		select {
		case <-s.stopped:
			break
		default:
			if !settled && s.Time() < 200 {
				// At the beginning we run the sequential simulator.
				s.Handle()
			} else {
				settled = true
				handler()
			}
		}
	}
}

// Handling events synchronously.
func (s *Simulation) Handle() {
	if event := s.Pop(); event != nil {
		// The event gets dispached to observers
		for _, observer := range s.observers {
			observer.Receive(event)
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

// Handling events in parallel.
//
// We allow parallel exectution following the next pipeline:
// - pop all events with the time-stamp equal with now
// - group all these events in event groups by the receiver address
// - execute the receivers in parallel
// - wait for all receivers to finish their execution
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
	for ; event != nil && event.timestamp == eventTime; event = s.Pop() {
		events = append(events, event)
	}
	// Push back the next event extracted.
	if event != nil {
		s.Push(event)
	}

	// The event gets dispached to observers.
	for _, event := range events {
		for _, observer := range s.observers {
			observer.Receive(event)
		}
	}

	// Modify the global timestamp
	s.timeMutex.Lock()
	s.time = events[0].timestamp
	s.timeMutex.Unlock()

	// Find special events.
	special := []*Event{}
	rest := []*Event{}
	for _, event := range events {
		if event == nil || event.Receiver() == nil {
			continue
		}
		if s.check(event.Receiver()) {
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
	for _, event := range events {
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
		done := make(chan bool)
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
