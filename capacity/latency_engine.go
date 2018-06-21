package capacity

import (
  "github.com/danalex97/Speer/overlay"
)

// A Transport which takes into consideration latency when sending control
// messages.
type TransferLatencyEngine struct {
  *TransferEngine

  unreliableNode overlay.UnreliableNode
}

func NewTransferLatencyEngine(
    e *TransferEngine,
    u overlay.UnreliableNode) Engine {

  engine := &TransferLatencyEngine{
    TransferEngine : e,
    unreliableNode : u,
  }
  engine.unreliableNode.SetProxy(engine.stripPayload)

  return engine
}

func (e *TransferLatencyEngine) ControlSend(id string, message interface{}) {
  e.unreliableNode.Send(overlay.NewPacket(
    e.id,
    id,
    message,
  ))
}

func (e *TransferLatencyEngine) stripPayload(m interface {}) interface {} {
  return m.(overlay.Packet).Payload()
}

func (e *TransferLatencyEngine) ControlRecv() <-chan interface{} {
  return e.unreliableNode.Recv()
}
