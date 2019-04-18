package interfaces

// The Util interface is provided to a node.
type TorrentNodeUtil interface {
	Transport() Transport

	Id() string
	Join() string

	Time() func() int
}

// This interface needs to be implemented by a node.
type TorrentNode interface {
	// The general method that is just a runner.
	OnJoin()

	// A method that should be called when a node leaves the network.
	OnLeave()
}

// This interface needs to be implemented by a node.
type TorrentNodeConstructor interface {
	New(util TorrentNodeUtil) TorrentNode
}
