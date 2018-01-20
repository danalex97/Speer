package main

import (
  "fmt"
  . "./events"
)

func main() {
  eq := NewHeapEventQueue()
  eq.Push(NewEvent(3, 10, nil))
  eq.Push(NewEvent(5, 10, nil))
  eq.Push(NewEvent(2, 10, nil))
  top := eq.Pop()
  fmt.Println(top.Timestamp())
  top = eq.Pop()
  fmt.Println(top.Timestamp())
  top = eq.Pop()
  fmt.Println(top.Timestamp())
}
