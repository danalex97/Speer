package capacity

import (
  . "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/Speer/overlay"
  "runtime"
  "fmt"
)

/* Implementation. */
type TransferLatencyEngine struct {
  *TransferEngine

  unreliableNode UnreliableNode
}

func NewTransferLatencyEngine(e *TransferEngine, u UnreliableNode) Engine {
  engine := &TransferLatencyEngine{
    TransferEngine : e,
    unreliableNode : u,
  }
  go engine.establishListener();
  return engine
}

func (e *TransferLatencyEngine) ControlSend(id string, message interface{}) {
  e.unreliableNode.Send() <- overlay.NewPacket(
    e.id,
    id,
    message,
  )
}

func (e *TransferLatencyEngine) establishListener() {
  for {
    select {
    case pkt := <- e.unreliableNode.Recv():
      if len(e.recv) == cap(e.recv) {
        fmt.Println("Channel blocked at ControlRecv.")
      }
      e.recv <- pkt.(overlay.Packet).Payload
    default:
      // If there are no new packets schedule other routines.
      runtime.Gosched()
    }
  }
}

func (e *TransferLatencyEngine) ControlRecv() <-chan interface{} {
  return e.recv
}
