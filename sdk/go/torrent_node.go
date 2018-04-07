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

  // used to autowire the engine
  autowireEngine(capacity.Engine)
}

type AutowiredTorrentNode struct {
  node     interfaces.UnreliableNode
  template interfaces.TorrentNode

  engine   capacity.Engine
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

/* Constructor */
func NewAutowiredTorrentNode(node interfaces.UnreliableNode, template interface {}) DHTNode {
  s := new(AutowiredTorrentNode)

  s.node     = node
  s.template = template.(interfaces.TorrentNodeConstructor).New(s)

  return s
}

/* Autowire. */
func (a *AutowiredTorrentNode) autowireEngine(engine capacity.Engine) {
  a.engine = engine
}
