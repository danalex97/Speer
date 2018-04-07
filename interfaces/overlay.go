package interfaces

import (
  "github.com/danalex97/Speer/overlay"
)

type UnreliableNode interface {
  Join()   string
  Id()     string

  Send()   chan<- interface{}
  Recv()   <-chan interface{}
}
