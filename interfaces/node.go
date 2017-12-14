package interfaces

// The Util interface is provided to a node.
type NodeUtil interface {
	Transport() Transport

	Id() string
	Join() string

	Time() func() int
}

// This interface needs to be implemented by a node.
type Node interface {
	// The general method that is just a runner.
	OnJoin()

	// A method that should be called when a node leaves the network.
	OnLeave()
}

// This interface needs to be implemented by a node.
type NodeConstructor interface {
	New(util NodeUtil) Node
}
