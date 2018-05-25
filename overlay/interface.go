package overlay

import (
  "github.com/danalex97/Speer/underlay"
  "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/Speer/events"
)

type UnreliableNode interface {
  events.Decorable

  interfaces.UnreliableNode
}

type UnreliableSimulatedNode struct {
  events.Decorable

  simulation *underlay.NetworkSimulation
  bridge     Bridge
  bootstrap  Bootstrap
  id         string
}

var activeSet = make(map[*underlay.NetworkSimulation]OverlayMap)

func GetBootstrap(simulation *underlay.NetworkSimulation) Bootstrap {
  netMap := NewNetworkMap(simulation.Network())
  if mp, ok := activeSet[simulation]; ok {
    netMap = mp
  } else {
    activeSet[simulation] = netMap
  }

  return netMap
}

func NewUnreliableSimulatedNode(simulation *underlay.NetworkSimulation) UnreliableNode {
  node := new(UnreliableSimulatedNode)

  var netMap OverlayMap
  if mp, ok := activeSet[simulation]; ok {
    netMap = mp
  } else {
    netMap = NewNetworkMap(simulation.Network())
    activeSet[simulation] = netMap
  }

  node.id         = netMap.NewId()
  node.bridge     = NewUnderlayChan(node.id, simulation, netMap)
  node.bootstrap  = netMap
  node.simulation = simulation

  // The actual decorable is at bridge level.
  // To allow direct interfacing, we create a tunnnel.
  node.Decorable = events.NewTunnel(node.bridge)

  return node
}

func (n *UnreliableSimulatedNode) Id() string {
  return n.id
}

func (n *UnreliableSimulatedNode) Send(msg interface {}) {
  n.bridge.Send(msg)
}

func (n *UnreliableSimulatedNode) Recv() <-chan interface{} {
  return n.bridge.Recv()
}

func (n *UnreliableSimulatedNode) Join() string {
  return n.bootstrap.Join(n.id)
}

func (n *UnreliableSimulatedNode) Simulation() *underlay.NetworkSimulation {
  return n.simulation
}
