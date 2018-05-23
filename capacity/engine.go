package capacity

import (
  . "github.com/danalex97/Speer/interfaces"
  "sync"
  "fmt"
)

/* Interface. */
type Engine interface {
  Transport

  // seam used to register the links in the SDK
  SetConnectCallback(func (Link))
}

/* Global variables */
var engineMap = make(map[string]Engine)
var mapLock   = new(sync.RWMutex)

/* There is only one control message queue, so it should be big. */
const controlMessageCapacity int = 1000000

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

  recv             chan interface {}
  id               string
  connectCallback  func (Link)
}

func NewTransferEngine(up, down int, id string) Engine {
  mapLock.Lock()
  defer mapLock.Unlock()

  engine := &TransferEngine{
    node{
      down : down,
      up   : up,
    },
    make(chan interface {}, controlMessageCapacity),
    id,
    func (Link) {},
  }

  engineMap[engine.id] = engine

  return engine
}

func (e *TransferEngine) Connect(id string) Link {
  mapLock.RLock()
  defer mapLock.RUnlock()

  link := NewPerfectLink(e, engineMap[id])
  e.connectCallback(link)

  return link
}

func (e *TransferEngine) ControlSend(id string, message interface{}) {
  mapLock.RLock()
  engine := engineMap[id].(*TransferEngine)
  mapLock.RUnlock()

  if len(engine.recv) == cap(engine.recv) {
    fmt.Println("Channel blocked at ControlSend.")
  }
  engine.recv <- message
}

func (e *TransferEngine) ControlRecv() <-chan interface{} {
  return e.recv
}

func (e *TransferEngine) ControlPing(id string) bool {
  mapLock.RLock()
  defer mapLock.RUnlock()

  return engineMap[id] != nil
}

func (e *TransferEngine) Id() string {
  return e.id
}

func (e *TransferEngine) SetConnectCallback(callback func (Link)) {
  e.connectCallback = callback
}
