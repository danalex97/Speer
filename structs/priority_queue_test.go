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

	assertEqual(t, q.Len(), 0)

	q.Push(Int(3), nil)
	q.Push(Int(1), nil)
	q.Push(Int(2), nil)

	assertEqual(t, q.Len(), 3)

	assertEqual(t, q.Pop().Key, Int(1))
	assertEqual(t, q.Pop().Key, Int(2))
	assertEqual(t, q.Pop().Key, Int(3))

	assertEqual(t, q.Len(), 0)
}
