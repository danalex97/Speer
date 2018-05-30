package capacity

import (
  . "github.com/danalex97/Speer/interfaces"
  "testing"
)

func TestCanConnectEngines(t *testing.T) {
  e1 := NewTransferEngine(10, 20, "1")
  NewTransferEngine(30, 40, "2")

  callbackCalled := false
  e1.SetConnectCallback(func (Link) {
    callbackCalled = true
  })

  l := e1.Connect("2")
  assertEqual(t, l.From().Up(), 10)
  assertEqual(t, l.From().Down(), 20)
  assertEqual(t, l.To().Up(), 30)
  assertEqual(t, l.To().Down(), 40)

  assertEqual(t, callbackCalled, true)
}

// func TestCanSendControlMessages(t *testing.T) {
//   e1 := NewTransferEngine(10, 20, "1")
//   e2 := NewTransferEngine(10, 20, "2")
//
//   e1.ControlSend("2", "message")
//   assertEqual(t, (<-e2.ControlRecv()).(string), "message")
// }
