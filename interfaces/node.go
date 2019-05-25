package interfaces

// The Util interface is provided to a node.
type NodeUtil interface {
	RoutineCapabilities

	Transport() Transport

	Id() string
	Join() string

	Time() func() int
}

// This interface needs to be implemented by a node. All the functions OnJoin,
// OnNotify and OnLeave should be non-blocking since the simulator uses
// coroutines.
type Node interface {
	// constructor interface
	New(util NodeUtil) Node

	// Function that represents what the node should do when notified.
	OnNotify()

	// Function that represent the initial action taken by a node when it
	// joins the network.
	OnJoin()

	// A method that should be called when a node leaves the network.
	OnLeave()
}
