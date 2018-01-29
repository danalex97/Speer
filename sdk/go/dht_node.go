package sdk

import (
  . "github.com/danalex97/Speer/overlay"
  . "github.com/danalex97/Speer/model"
  "github.com/danalex97/Speer/underlay"
)

type DHTNode interface {
  UnreliableNode() UnreliableNode
  // an unreliable node interface is a mean of interaction with an
  // underlay simulation through the overlay

  OnJoin()
  // a method that should be called when a node joins the network

  OnQuery(query DHTQuery) error
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
}

// autowiring mechanism to hide simulation injection at construction
type Autowire interface {
  DHTNode
  autowire() Autowire
}

type AutowiredDHTNode struct {
  node UnreliableNode
}

func (a *AutowiredDHTNode) autowire() Autowire {
  return a
}

func (a *AutowiredDHTNode) UnreliableNode() UnreliableNode {
  return s.node;
}
