package interfaces

type UnreliableNode interface {
  Join()   string
  Id()     string

  Send()   chan<- interface{}
  Recv()   <-chan interface{}
}
