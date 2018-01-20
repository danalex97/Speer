package main

import (
  . "github.com/danalex97/Speer/events"
  . "github.com/danalex97/Speer/underlay"
)

func main() {
  network := NewRandomUniformNetwork(10, 40, 1, 10)
  packet1 := NewPacket(network.Routers[0], network.Routers[5])
  packet2 := NewPacket(network.Routers[1], network.Routers[2])

  s := NewNetworkSimulation(NewLazySimulation(), network)

  s.SendPacket(packet1)
  s.SendPacket(packet2)

  s.Run()
  s.Stop()
}
