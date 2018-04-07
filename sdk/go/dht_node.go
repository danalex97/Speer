package sdk

import (
  "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/Speer/overlay"
)

type DHTNode interface {
  UnreliableNode() overlay.UnreliableNode
  // an unreliable node interface is a mean of interaction with an
  // underlay simulation through the overlay

  OnJoin()
  // a method that should be called when a node joins the network

  OnQuery(query interfaces.Query) error
  // a method that should be called with a node receives a query
  // - the query can be either:
  //   - store -- the current node *wants* to store something in the network
  //   - query -- the current node *wants* to retrieve something from the network
  // both operations could result in an error

  OnLeave()
  // a meothd that should be called when a node leaves the network

  NewDHTNode() DHTNode
  // used to generate a DHTNode

  Key() string
  // generate a new key for the key space

  Autowire(template DHTNode)
  // used to autowire the node to the simulation

  autowire() Autowire
  // used to autowire the node to the simulation
}

func autowiredUnreliableNode(node DHTNode) overlay.UnreliableNode {
  simulation := node.autowire().(*AutowiredDHTNode).node.(*overlay.UnreliableSimulatedNode).Simulation()
  newNode := overlay.NewUnreliableSimulatedNode(simulation)

  return newNode
}

// autowiring mechanism to hide simulation injection at construction
type Autowire interface {
  UnreliableNode() overlay.UnreliableNode
  autowire() Autowire
  Autowire(template DHTNode)
}

type AutowiredDHTNode struct {
  node overlay.UnreliableNode
}

func (a *AutowiredDHTNode) Autowire(template DHTNode) {
  unreliableNode := autowiredUnreliableNode(template)
  a.autowire().(*AutowiredDHTNode).node = unreliableNode
}

func (a *AutowiredDHTNode) autowire() Autowire {
  return a
}

func (a *AutowiredDHTNode) UnreliableNode() overlay.UnreliableNode {
  return a.node
}
