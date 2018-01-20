package main

import (
  . "github.com/danalex97/Speer/events"
  . "github.com/danalex97/Speer/underlay"
  "math/rand"
  "time"
)

func main() {
  rand.Seed(time.Now().UTC().UnixNano())

  network := NewRandomUniformNetwork(10000, 70000, 2, 10)
  packet1 := NewPacket(network.Routers[0], network.Routers[5])
  packet2 := NewPacket(network.Routers[1], network.Routers[2])

  s := NewNetworkSimulation(NewLazySimulation(), network)

  s.SendPacket(packet1)
  s.SendPacket(packet2)

  go s.Run()

  time.Sleep(time.Duration(1) * time.Second)
  s.Stop()
}
