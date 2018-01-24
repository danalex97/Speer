package main

import (
  "github.com/danalex97/Speer/events"
  "github.com/danalex97/Speer/underlay"
  "github.com/danalex97/Speer/overlay"
  "math/rand"
  "time"
  "fmt"
)

func main() {
  rand.Seed(time.Now().UTC().UnixNano())

  network := underlay.NewRandomUniformNetwork(10000, 70000, 2, 10)
  s := underlay.NewNetworkSimulation(events.NewLazySimulation(), network)

  node1 := overlay.NewUnreliableSimulatedNode(s)
  node2 := overlay.NewUnreliableSimulatedNode(s)
  node1.Send() <- overlay.NewPacket(node1.Id(), node2.Id(), nil)

  go s.Run()
  time.Sleep(time.Duration(1) * time.Second)
  s.Stop()

  fmt.Println("Stopped simulation")
}
