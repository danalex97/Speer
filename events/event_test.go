package events

import (
  "testing"
)

func TestNewEvent(t *testing.T) {
  r := new(mockReceiver)
  e := NewEvent(10, nil, r)

  assertEqual(t, e.Timestamp(), 10)
  assertEqual(t, e.Payload(), nil)
  assertEqual(t, e.Receiver(), r)
}
