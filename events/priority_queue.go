package events

import (
  "container/heap"
)

type EventQueue interface {
  Push(e *Event)
  Pop() *Event
}

type eventHeap []*Event

func (h eventHeap) Len() int           { return len(h) }
func (h eventHeap) Less(i, j int) bool { return h[i].timestamp < h[j].timestamp }
func (h eventHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *eventHeap) Push(x interface{}) {
	*h = append(*h, x.(*Event))
}

func (h *eventHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type heapEventQueue struct {
  h heap.Interface
}

func NewHeapEventQueue() EventQueue {
  h := &eventHeap{}
  heap.Init(h)

  eq  := new(heapEventQueue)
  eq.h = h
  return eq
}

func (eq *heapEventQueue) Push(e *Event) {
  heap.Push(eq.h, e)
}

func (eq *heapEventQueue) Pop() *Event {
  return heap.Pop(eq.h).(*Event)
}
