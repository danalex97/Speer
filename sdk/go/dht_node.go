package sdk

import (
	"github.com/danalex97/Speer/interfaces"
)

type DHTNode interface {
	interfaces.DHTNode
	interfaces.DHTNodeUtil
}

// An AutowiredDHTNode is a DHT node which can be used in a simulation. The
// user needs to provide an interfaces.DHTNode template for this purpose.
type AutowiredDHTNode struct {
	node     interfaces.UnreliableNode
	template interfaces.DHTNode
}

// Part of DHTNode interface.
func (n *AutowiredDHTNode) OnJoin() {
	n.template.OnJoin()
}

// Part of DHTNode interface.
func (n *AutowiredDHTNode) OnQuery(query interfaces.Query) error {
	return n.template.OnQuery(query)
}

// Part of DHTNode interface.
func (n *AutowiredDHTNode) OnLeave() {
	n.template.OnLeave()
}

// Part of DHTNode interface.
func (n *AutowiredDHTNode) Key() string {
	return n.template.Key()
}

// Part of DHTNodeUtil interface.
func (n *AutowiredDHTNode) UnreliableNode() interfaces.UnreliableNode {
	return n.node
}

// Contructor and decorator.
func NewAutowiredDHTNode(node interfaces.UnreliableNode, simulation interface{}) DHTNode {
	n := new(AutowiredDHTNode)

	s := simulation.(*DHTSimulation)

	n.node = node
	n.template = s.template.(interfaces.DHTNodeConstructor).New(n)

	return n
}
