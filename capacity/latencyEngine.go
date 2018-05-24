package capacity

import (
  . "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/Speer/overlay"
  "fmt"
)

type TransmissionProgress *overlay.TransmissionProgress

/* Implementation. */
type TransferLatencyEngine struct {
  *TransferEngine

  unreliableNode UnreliableNode
  prog           TransmissionProgress
}

func NewTransferLatencyEngine(
    e *TransferEngine,
    u UnreliableNode,
    p TransmissionProgress) Engine {

  engine := &TransferLatencyEngine{
    TransferEngine : e,
    unreliableNode : u,
    prog           : p,
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
    pkt := <-e.unreliableNode.Recv()
    if len(e.recv) == cap(e.recv) {
      fmt.Println("Channel blocked at ControlRecv.")
    }
    e.recv <- pkt.(overlay.Packet).Payload()
  }
}

func (e *TransferLatencyEngine) ControlRecv() <-chan interface{} {
  return e.recv
}
