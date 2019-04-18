// package capacity
//
// import (
//   . "github.com/danalex97/Speer/interfaces"
//
//   "github.com/danalex97/Speer/overlay"
//   "github.com/danalex97/Speer/events"
//   "github.com/danalex97/Speer/structs"
//
//   "testing"
// )
//
// type unreliableNode struct {
//   *events.Decorator
//
//   chn chan interface {}
// }
//
// func newMockUnreliableNode() overlay.UnreliableNode {
//   return &unreliableNode{
//     Decorator : events.NewDecorator(),
//     chn : make(chan interface {}, 100),
//   }
// }
//
// func (u *unreliableNode) Join() string {
//   return structs.RandomKey()
// }
//
// func (u *unreliableNode) Id() string {
//   return structs.RandomKey()
// }
//
// func (u *unreliableNode) Send(v interface {}) {
//   u.chn <- u.Proxy(v)
// }
//
// func (u *unreliableNode) Recv() <-chan interface {} {
//   return u.chn
// }
//
// func TestCanConnectLatencyEngines(t *testing.T) {
//   e1 := NewTransferLatencyEngine(
//     NewTransferEngine(10, 20, "1").(*TransferEngine),
//     newMockUnreliableNode(),
//   )
//   NewTransferEngine(30, 40, "2")
//
//   callbackCalled := false
//   e1.SetConnectCallback(func (Link) {
//     callbackCalled = true
//   })
//
//   l := e1.Connect("2")
//   assertEqual(t, l.From().Up(), 10)
//   assertEqual(t, l.From().Down(), 20)
//   assertEqual(t, l.To().Up(), 30)
//   assertEqual(t, l.To().Down(), 40)
//
//   assertEqual(t, callbackCalled, true)
// }
//
// func TestCanSendControlMessages(t *testing.T) {
//   e := NewTransferLatencyEngine(
//     NewTransferEngine(10, 20, "1").(*TransferEngine),
//     newMockUnreliableNode(),
//   )
//   e.ControlSend("1", "hmm")
//   assertEqual(t, (<-e.ControlRecv()).(string), "hmm")
// }
