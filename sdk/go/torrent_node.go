package sdk

import (
  "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/Speer/capacity"
  "github.com/danalex97/Speer/structs"
)

type TorrentNode interface {
  interfaces.DHTNode
  interfaces.DHTNodeUtil
  interfaces.TorrentNodeUtil
}

type AutowiredTorrentNode struct {
  node     interfaces.UnreliableNode
  template interfaces.TorrentNode

  engine   capacity.Engine
  time     func() int
}

// Part of DHTNode interface.
func (n *AutowiredTorrentNode) OnJoin() {
  n.template.OnJoin()
}

// Part of DHTNode interface.
func (n *AutowiredTorrentNode) OnQuery(query interfaces.Query) error {
  return nil
}

// Part of DHTNode interface.
func (n *AutowiredTorrentNode) OnLeave() {
  n.template.OnLeave()
}

// Part of DHTNode interface.
func (n *AutowiredTorrentNode) Key() string {
  return structs.RandomKey()
}

// Part of DHTNodeUtil interface.
func (n *AutowiredTorrentNode) UnreliableNode() interfaces.UnreliableNode {
  return n.node
}

// Part of TorrentNodeUtil interface.
func (n *AutowiredTorrentNode) Transport() interfaces.Transport {
  return n.engine
}

// Part of TorrentNodeUtil interface.
func (n *AutowiredTorrentNode) Id() string {
  return n.node.Id()
}

// Part of TorrentNodeUtil interface.
func (n *AutowiredTorrentNode) Join() string {
  return n.node.Join()
}

// Part of TorrentNodeUtil interface.
func (n *AutowiredTorrentNode) Time() func() int {
  return n.time
}

// Constructor.
func NewAutowiredTorrentNode(node interfaces.UnreliableNode, simulation interface {}) DHTNode {
  n := new(AutowiredTorrentNode)

  s := simulation.(*TorrentSimulation)

  n.node     = node
  n.engine   = s.updateEngine(node)
  n.time     = s.Time
  n.template = s.template.(interfaces.TorrentNodeConstructor).New(n)

  return n
}
