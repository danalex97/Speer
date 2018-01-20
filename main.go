package main

import (
  . "github.com/danalex97/Speer/events"
  . "github.com/danalex97/Speer/underlay"
)

func main() {
  s := NewLazySimulation()
  go s.Run()

  go s.Push(NewEvent(3, 10, nil))
  go s.Push(NewEvent(5, 10, nil))
  go s.Push(NewEvent(2, 10, nil))

  s.Stop()

  r1 := NewShortestPathRouter()
  r2 := NewShortestPathRouter()
  r3 := NewShortestPathRouter()

  r1.Connect(NewStaticConnection(5, r2))
  r2.Connect(NewStaticConnection(5, r3))
  r1.Receive(NewEvent(0, *NewPacket(r1, r3), r1))
}
