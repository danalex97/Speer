package events

import (
	"testing"
)

func TestRoutineCallsFunctionEachTimeAndCanBeStopped(t *testing.T) {
	v := 1
	r := NewRoutine(5, func() { v++ })
	e := NewEvent(0, nil, nil)
	for i := 1; i <= 5; i++ {
		assertEqual(t, v, i)
		e = r.Receive(e)
		assertEqual(t, e.Timestamp(), i * 5);
	}
	assertEqual(t, v, 6)
	r.Stop()
	e = r.Receive(e)
	assertEqual(t, v, 6)
	assertEqual(t, e, (*Event)(nil))
}

func TestCallbackIsCalledOnce(t *testing.T) {
	v := 1
	c := NewCallback(func() { v++ })
	e := NewEvent(0, nil, nil)
	e = c.Receive(e)
	assertEqual(t, v, 2)
	assertEqual(t, e, (*Event)(nil))
}

func TestCallbackCanBeStopped(t *testing.T) {
	v := 1
	c := NewCallback(func() { v++ })
	e := NewEvent(0, nil, nil)
	c.Stop()
	e = c.Receive(e)
	assertEqual(t, v, 1)
	assertEqual(t, e, (*Event)(nil))
}
