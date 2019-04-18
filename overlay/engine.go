package overlay

import (
  . "github.com/danalex97/Speer/capacity"
  "github.com/danalex97/Speer/interfaces"
  "sync"
  "fmt"
)

// The Engine is a Transport decorated with methods which allow the integration
// of the capacity engine in the SDK.
type Engine interface {
  interfaces.Transport

  // seam used to register the links in the SDK
  SetConnectCallback(func (interfaces.Link))
}

// Global variables
var engineMap = make(map[string]Engine)
var mapLock   = new(sync.RWMutex)

// There is only one control message queue, so it should be big.
const controlMessageCapacity int = 1000000

// Simple Node interface structure.
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

// A TransferEngine is a interfaces.Transport which allows direct transfer of control
// messages and uses links simulated using the Scheduler.
type TransferEngine struct {
  node

  recv             chan interface {}
  id               string
  connectCallback  func (interfaces.Link)
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
    func (interfaces.Link) {},
  }

  engineMap[engine.id] = engine

  return engine
}

func (e *TransferEngine) Connect(id string) interfaces.Link {
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

func (e *TransferEngine) SetConnectCallback(callback func (interfaces.Link)) {
  e.connectCallback = callback
}
