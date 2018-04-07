package interfaces

/* The Util interface is provided to a node. */
type DHTNodeUtil interface {
  UnreliableNode() UnreliableNode
  // an unreliable node interface is a mean of interaction with an
  // underlay simulation through the overlay
}

/* This interface has to be implemented by a node. */
type DHTNode interface {
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

  New(util DHTNodeUtil)
  // the constructor called by the simulation

  Key() string
  // generate a new key for the key space
}
