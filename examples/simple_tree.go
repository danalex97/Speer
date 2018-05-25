package examples

import (
  . "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/Speer/overlay"
  "github.com/danalex97/Speer/structs"
  "runtime"
  "sync"
  "fmt"
)

type SimpleTree struct {
  sync.Mutex

  id           string
  neighId      string
  store        map[string]bool

  unreliableNode UnreliableNode
}

func (s *SimpleTree) OnJoin() {
  go func() {
    for {
      select {
      case _, ok := <-s.unreliableNode.Recv():
        if ok {
          fmt.Println("Receive")
        }
      default:
        runtime.Gosched()
      }
    }
  }()
}

func (s *SimpleTree) OnQuery(query Query) error {
  s.Lock()
  defer s.Unlock()

  key := query.Key()
  if query.Store() {
    key = s.Key()
    s.store[key] = true
  } else {
    // check in my local store
    if _, ok := s.store[key]; ok {
      return nil
    }

    // check the other node's store to retrieve
    packet := overlay.NewPacket(
      s.id,
      s.neighId,
      query,
    )
    s.unreliableNode.Send(packet)
  }

  return nil
}

func (s *SimpleTree) OnLeave() {
}

func (s *SimpleTree) New(util DHTNodeUtil) DHTNode {
  // Constructor that assumes the UnreliableNode component is filled in
  node := new(SimpleTree)

  node.id       = util.UnreliableNode().Id()
  node.neighId  = util.UnreliableNode().Join()
  node.unreliableNode = util.UnreliableNode()
  node.store    = make(map[string]bool)

  return node
}

func (s *SimpleTree) Key() string {
  return structs.RandomKey()
}
