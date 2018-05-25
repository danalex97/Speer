package interfaces

type UnreliableNode interface {
  Join() string
  Id()   string

  Send(interface {})
  Recv() <-chan interface{}
}
