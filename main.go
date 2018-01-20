package main

import (
  . "github.com/danalex97/Speer/events"
  . "github.com/danalex97/Speer/underlay"
)

func main() {
  s := NewLazySimulation()
  go s.Run()

  s.Push(NewEvent(3, 10, nil))
  s.Push(NewEvent(5, 10, nil))
  s.Push(NewEvent(2, 10, nil))

  network := NewRandomUniformNetwork(10, 40, 1, 5)
  s.Push(NewEvent(1, NewPacket(network.Routers[0], network.Routers[5]), network.Routers[0]))
  s.Push(NewEvent(5, NewPacket(network.Routers[3], network.Routers[2]), network.Routers[3]))

  s.Stop()
}
