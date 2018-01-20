package main

import (
  . "./events"
)

func main() {
  s := NewLazySimulation()
  go s.Run()

  go s.Push(NewEvent(3, 10, nil))
  go s.Push(NewEvent(5, 10, nil))
  go s.Push(NewEvent(2, 10, nil))

  s.Stop()
}
