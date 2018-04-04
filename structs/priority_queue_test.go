package structs

import (
  "testing"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

func TestPriorityQueueOrderedElements(t *testing.T) {
  q := NewPriorityQueue()

  q.Push(3, nil)
  q.Push(1, nil)
  q.Push(2, nil)

  assertEqual(t, q.Pop().Key, 1)
  assertEqual(t, q.Pop().Key, 2)
  assertEqual(t, q.Pop().Key, 3)
}
