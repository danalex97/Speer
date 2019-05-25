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
	i int
}

func (m *mockReceiver) Receive(e *Event) *Event {
	return nil
}

func TestProxyIsCalledByActiveObserver(t *testing.T) {
	r1 := &mockReceiver{}
	r2 := &mockReceiver{}
	o := NewActiveEventObserver(r1)

	ctr := 0
	o.SetProxy(NewProxy(func () {
		ctr += 1
	}))

	e1 := NewEvent(0, nil, r1)
	e2 := NewEvent(0, nil, r2)
	o.Receive(e1)
	o.Receive(e2)

	assertEqual(t, ctr, 1)
}

func TestEventsAreDeliveredToRightPassiveObserver(t *testing.T) {
	r1 := new(mockReceiver)
	r2 := new(mockReceiver)

	o := NewPassiveEventObserver(r1)

	e1 := NewEvent(1, nil, r1)
	e2 := NewEvent(0, nil, r2)

	go func() {
		o.Receive(e1)
		o.Receive(e2)
	}()

	go func() {
		assertEqual(t, (<-o.Recv()).(*Event), e1)
	}()
}

func TestEnquedEventsAreDeliveredInOrder(t *testing.T) {
	r := new(mockReceiver)
	o := NewPassiveEventObserver(r)

	e1 := NewEvent(0, nil, r)
	e2 := NewEvent(1, nil, r)
	e3 := NewEvent(2, nil, r)

	go func() {
		o.Receive(e1)
		o.Receive(e2)
		o.Receive(e3)
	}()

	go func() {
		assertEqual(t, (<-o.Recv()).(*Event), e1)
		assertEqual(t, (<-o.Recv()).(*Event), e2)
		assertEqual(t, (<-o.Recv()).(*Event), e3)
	}()
}

func TestAllEventsAreDeliveredToGlobalObserver(t *testing.T) {
	r1 := new(mockReceiver)
	r2 := new(mockReceiver)
	r3 := new(mockReceiver)
	o := NewGlobalEventObserver()

	e1 := NewEvent(0, nil, r1)
	e2 := NewEvent(1, nil, r2)
	e3 := NewEvent(2, nil, r3)

	go func() {
		o.Receive(e1)
		o.Receive(e2)
		o.Receive(e3)
	}()

	go func() {
		assertEqual(t, (<-o.Recv()).(*Event), e1)
		assertEqual(t, (<-o.Recv()).(*Event), e2)
		assertEqual(t, (<-o.Recv()).(*Event), e3)
	}()
}
