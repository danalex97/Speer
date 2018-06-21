package interfaces

// The Util interface is provided to a node.
type DHTNodeUtil interface {
  // an unreliable node interface is a mean of interaction with an
  // underlay simulation through the overlay
  UnreliableNode() UnreliableNode
}

// This interface has to be implemented by a node.
type DHTNode interface {
  // OnJoin is a method that should be called when a node joins the network.
  OnJoin()

  // OnQuery is a method that should be called with a node receives a query.
  // The query can be either:
  //  - store -- the current node *wants* to store something in the network
  //  - query -- the current node *wants* to retrieve something from the network
  // Both operations could result in an error.
  OnQuery(query Query) error

  // OnLeave is a method that should be called when a node leaves the network.
  OnLeave()

  // The Key function is used to generate a new key for the key space.
  Key() string
}

// This interface needs to be implemented by a node.
type DHTNodeConstructor interface {
  New(util DHTNodeUtil) DHTNode
}
