package underlay

import (
  . "github.com/danalex97/Speer/events"
)

type NetworkSimulation struct {
  *Simulation
  network *Network
}

func NewNetworkSimulation(s *Simulation, n *Network) *NetworkSimulation {
  return &NetworkSimulation{
    Simulation : s,
    network    : n,
  }
}

func (s *NetworkSimulation) SendPacket(p Packet) {
  s.Push(NewEvent(s.Time(), p, p.Src()))
}

func (s *NetworkSimulation) Network() *Network {
  return s.network
}
