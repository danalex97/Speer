package events

import (
  "testing"
)

func TestEventQueueDeliversOrderedElements(t *testing.T) {
  r := new(mockReceiver)

  e := NewLazyEventQueue()

  e1 := NewEvent(0, nil, r)
  e2 := NewEvent(1, nil, r)
  e3 := NewEvent(2, nil, r)

  e.Push(e3)
  e.Push(e1)
  e.Push(e2)

  assertEqual(t, e1, e.Pop())
  assertEqual(t, e2, e.Pop())
  assertEqual(t, e3, e.Pop())
}

func TestEventQueueDeliversOrderedElementsDifferentPushers(t *testing.T) {
  r := new(mockReceiver)

  e := NewLazyEventQueue()

  e1 := NewEvent(0, nil, r)
  e2 := NewEvent(1, nil, r)
  e3 := NewEvent(2, nil, r)

  done := make(chan bool)

  go func() {
    e.Push(e3)
    done <- true
  }()
  go func() {
    e.Push(e1)
    done <- true
  }()
  go func() {
    e.Push(e2)
    done <- true
  }()

  for i := 0; i < 3; i++ {
    <-done
  }

  assertEqual(t, e1, e.Pop())
  assertEqual(t, e2, e.Pop())
  assertEqual(t, e3, e.Pop())
}
