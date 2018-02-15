package events

import (
  "container/heap"
)

type EventQueue interface {
  Push(e *Event)
  Pop() *Event
}

const LazyQueueChanSize int = 50

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

type lazyEventQueue struct {
  h heap.Interface
  stream chan *Event
}

func NewLazyEventQueue() EventQueue {
  h := &eventHeap{}
  heap.Init(h)

  eq := new(lazyEventQueue)

  eq.h      = h
  eq.stream = make(chan *Event, LazyQueueChanSize + 5)
  return eq
}

func (eq *lazyEventQueue) Push(event *Event) {
  eq.stream <- event
}

func (eq *lazyEventQueue) Pop() (event *Event) {
  eq.stream <- nil
  for {
    event = <- eq.stream
    if event == nil {
      break
    }

    heap.Push(eq.h, event)
  }

  if eq.h.Len() > 0 {
    event = heap.Pop(eq.h).(*Event)
  }
  return
}
