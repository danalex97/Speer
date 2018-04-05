package capacity

type Engine interface {
  Node

  Connect(Engine) Link

  ControlSend(Engine, interface {})
  ControlRecv() <-chan interface {}
}

const controlMessageCapacity int = 50

type node struct {
  down int
  up   int
}

func (n *node) Up() int {
  return n.up
}

func (n *node) Down() int {
  return n.down
}

type TransferEngine struct {
  node
  recv chan interface {}
}

func NewTransferEngine(n Node) Engine {
  return &TransferEngine{
    node{
      n.Up(),
      n.Down(),
    },
    make(chan interface {}, controlMessageCapacity),
  }
}

func (e *TransferEngine) Connect(node Engine) Link {
  return NewPerfectLink(e, node)
}

func (e *TransferEngine) ControlSend(engine Engine, message interface{}) {
  engine.(*TransferEngine).recv <- message
}

func (e *TransferEngine) ControlRecv() <-chan interface{} {
  return e.recv
}
