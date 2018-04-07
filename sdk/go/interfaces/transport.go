package interfaces

import (
  "github.com/danalex97/Speer/capacity"
)

type Data capacity.Data
/**
type Data struct {
  Id   string
  Size int
}
*/

type Node capacity.Node
/**
type Node interface {
  Up()   int
  Down() int
}
*/

type Link capacity.Link
/**
type Link interface {
  Upload(Data)
  Download() <-chan Data

  From() Node
  To()   Node
}
**/

type Transport capacity.Engine
/**
type Transport interface {
  Node

  Connect(string) Link

  ControlPing(string) bool
  ControlSend(string, interface {})
  ControlRecv() <-chan interface {}
}
*/
