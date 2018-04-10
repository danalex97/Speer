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

/* DHTNode interface */
func (n *AutowiredTorrentNode) OnJoin() {
  n.template.OnJoin()
}

func (n *AutowiredTorrentNode) OnQuery(query interfaces.Query) error {
  return nil
}

func (n *AutowiredTorrentNode) OnLeave() {
  n.template.OnLeave()
}

func (n *AutowiredTorrentNode) Key() string {
  return structs.RandomKey()
}

/* DHTNodeUtil interface */
func (n *AutowiredTorrentNode) UnreliableNode() interfaces.UnreliableNode {
  return n.node
}

/* TorrentNodeUtil interface */
func (n *AutowiredTorrentNode) Transport() interfaces.Transport {
  return n.engine
}

func (n *AutowiredTorrentNode) Id() string {
  return n.node.Id()
}

func (n *AutowiredTorrentNode) Join() string {
  return n.node.Join()
}

func (n *AutowiredTorrentNode) Time() func() int {
  return n.time
}

/* Constructor */
func NewAutowiredTorrentNode(node interfaces.UnreliableNode, simulation interface {}) DHTNode {
  n := new(AutowiredTorrentNode)

  s := simulation.(*TorrentSimulation)

  n.node     = node
  n.engine   = s.updateEngine(node)
  n.time     = s.Time
  n.template = s.template.(interfaces.TorrentNodeConstructor).New(n)

  return n
}
