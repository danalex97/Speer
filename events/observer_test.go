package events

import (
  "testing"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

type mockReceiver struct {
}

func (m *mockReceiver) Receive(e *Event) *Event {
  return nil
}

func TestEventsAreDeliveredToRightObserver(t *testing.T) {
  r1  := new(mockReceiver)
  r2  := new(mockReceiver)

  o := NewEventObserver(r1)

  e1 := NewEvent(1, nil, r1)
  e2 := NewEvent(0, nil, r2)

  go func() {
    o.EnqueEvent(e1)
    o.EnqueEvent(e2)
  }()

  go func() {
    assertEqual(t, <-o.EventChan(), e1)
  }()
}

func TestEnquedEventsAreDeliveredInOrder(t *testing.T) {
  r  := new(mockReceiver)
  o := NewEventObserver(r)

  e1 := NewEvent(0, nil, r)
  e2 := NewEvent(1, nil, r)
  e3 := NewEvent(2, nil, r)

  go func() {
    o.EnqueEvent(e1)
    o.EnqueEvent(e2)
    o.EnqueEvent(e3)
  }()

  go func() {
    assertEqual(t, <-o.EventChan(), e1)
    assertEqual(t, <-o.EventChan(), e2)
    assertEqual(t, <-o.EventChan(), e3)
  }()
}
