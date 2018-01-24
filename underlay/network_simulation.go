package underlay

import (
  . "github.com/danalex97/Speer/events"
)

type NetworkSimulation struct {
  Simulation
  network *Network
}

func NewNetworkSimulation(s Simulation, n *Network) *NetworkSimulation {
  ns := new(NetworkSimulation)
  ns.Simulation = s
  ns.network = n
  return ns
}

func (s *NetworkSimulation) SendPacket(p Packet) {
  s.Push(NewEvent(s.Time(), p, p.Src()))
}

func (s *NetworkSimulation) Network() *Network {
  return s.network
}
