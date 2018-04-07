package interfaces

type TorrentNodeUtil interface {
  Transport() Transport
  
  Id()   string
  Join() string
}

/* This interface needs to be implemented by a node.*/
type TorrentNode interface {
  OnJoin()
  // the general method that is just a runner

  OnLeave()
  // a method that should be called when a node leaves the network
}

/* This interface needs to be implemented by a node.*/
type TorrentNodeConstructor interface {
  New(util TorrentNodeUtil) TorrentNode
}
