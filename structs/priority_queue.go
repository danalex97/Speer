package structs

import (
  "container/heap"
)

type Element struct {
  Key   Comparable
  Value interface {}
}

type PriorityQueue interface {
  Len() int
  Push(Comparable, interface {})
  Pop() *Element
}

type container []*Element

func (c container) Len() int           { return len(c) }
func (c container) Less(i, j int) bool { return c[i].Key.Less(c[j].Key) }
func (c container) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

func (c *container) Push(x interface{}) {
  *c = append(*c, x.(*Element))
}

func (c *container) Pop() interface {} {
  old := *c
  n   := len(old)

  x := old[n - 1]
  *c =  old[0 : n - 1]

  return x
}

type pq struct {
  h heap.Interface
}

func NewPriorityQueue() PriorityQueue {
  h := &container{}
  heap.Init(h)

  q := new(pq)
  q.h = h

  return q
}

func (q *pq) Push(key Comparable, value interface {}) {
  heap.Push(q.h, &Element{key, value})
}

func (q *pq) Pop() *Element {
  return heap.Pop(q.h).(*Element)
}

func (q *pq) Len() int {
  return q.h.Len()
}
