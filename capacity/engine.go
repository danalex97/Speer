package capacity

import (
  "sync"
)

/* Interface. */
type Engine interface {
  Node

  Id() string
  Connect(string) Link

  ControlSend(string, interface {})
  ControlRecv() <-chan interface {}
}

/* Global variables */
var engineMap = *new(map[string]Engine)
var ctr       = 0
var mapLock   = new(sync.RWMutex)

const controlMessageCapacity int = 50

/* Simple node interface structure. */
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

/* Implementation. */
type TransferEngine struct {
  node
  recv chan interface {}
  id   string
}

func NewTransferEngine(up, down int) Engine {
  mapLock.Lock()
  defer mapLock.Unlock()

  engine := &TransferEngine{
    node{
      up,
      down,
    },
    make(chan interface {}, controlMessageCapacity),
    string(ctr),
  }

  engineMap[engine.id] = engine
  ctr++

  return engine
}

func (e *TransferEngine) Connect(id string) Link {
  mapLock.RLock()
  defer mapLock.RUnlock()

  return NewPerfectLink(e, engineMap[id])
}

func (e *TransferEngine) ControlSend(id string, message interface{}) {
  mapLock.RLock()
  defer mapLock.RUnlock()

  engineMap[id].(*TransferEngine).recv <- message
}

func (e *TransferEngine) ControlRecv() <-chan interface{} {
  return e.recv
}

func (e *TransferEngine) Id() string {
  return e.id
}
