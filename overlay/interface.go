package overlay

import (
  "github.com/danalex97/Speer/underlay"
)

type UnreliableNode interface {
  Id()   string
  Join() string
  Bridge
}

type UnreliableSimulatedNode struct {
  simulation *underlay.NetworkSimulation
  bridge     Bridge
  bootstrap  Bootstrap
  id         string
}

func NewUnreliableSimulatedNode(simulation *underlay.NetworkSimulation) UnreliableNode {
  node := new(UnreliableSimulatedNode)
  netMap := NewNetworkMap(simulation.Network())

  node.id         = netMap.NewId()
  node.bridge     = NewUnderlayChan(node.id, simulation, netMap)
  node.bootstrap  = netMap
  node.simulation = simulation

  return node
}

func (n *UnreliableSimulatedNode) Id() string {
  return n.id
}

func (n *UnreliableSimulatedNode) Send() chan<- interface{} {
  return n.bridge.Send()
}

func (n *UnreliableSimulatedNode) Recv() <-chan interface{} {
  return n.bridge.Recv()
}

func (n *UnreliableSimulatedNode) Join() string {
  return n.bootstrap.Join(n.id)
}
