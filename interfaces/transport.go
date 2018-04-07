package interfaces

type Data struct {
  Id   string
  Size int
}

type Node interface {
  Up()   int
  Down() int
}

type Link interface {
  Upload(Data)
  Download() <-chan Data

  From() Node
  To()   Node
}

type Transport interface {
  Node

  Connect(string) Link

  ControlPing(string) bool
  ControlSend(string, interface {})
  ControlRecv() <-chan interface {}
}
