package sdk

import (
  "github.com/danalex97/Speer/interfaces"
)

type DHTNode interface {
  interfaces.DHTNode
  interfaces.DHTNodeUtil
}

type AutowiredDHTNode struct {
  node     interfaces.UnreliableNode
  template interfaces.DHTNode
}

/* DHTNode interface */
func (n *AutowiredDHTNode) OnJoin() {
  n.template.OnJoin()
}

func (n *AutowiredDHTNode) OnQuery(query interfaces.Query) error {
  return n.template.OnQuery(query)
}

func (n *AutowiredDHTNode) OnLeave() {
  n.template.OnLeave()
}

func (n *AutowiredDHTNode) Key() string {
  return n.template.Key()
}

/* DHTNodeUtil interface */
func (n *AutowiredDHTNode) UnreliableNode() interfaces.UnreliableNode {
  return n.node
}

/* Constructor and decorator. */
func NewAutowiredDHTNode(node interfaces.UnreliableNode, simulation interface {}) DHTNode {
  n := new(AutowiredDHTNode)

  s := simulation.(*DHTSimulation)

  n.node     = node
  n.template = s.template.(interfaces.DHTNodeConstructor).New(n)

  return n
}
