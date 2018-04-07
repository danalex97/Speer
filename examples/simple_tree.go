package examples

import (
  . "github.com/danalex97/Speer/sdk/go"
  "github.com/danalex97/Speer/overlay"
  "github.com/danalex97/Speer/interfaces"
  "runtime"
  "sync"
  "fmt"
)

type SimpleTree struct {
  AutowiredDHTNode
  sync.Mutex

  id           string
  neighId      string
  store        map[string]bool
}

func (s *SimpleTree) OnJoin() {
  go func() {
    for {
      select {
      case _, ok := <-s.UnreliableNode().Recv():
        if ok {
          fmt.Println("Receive")
        }
      default:
        runtime.Gosched()
      }
    }
  }()
}

func (s *SimpleTree) OnQuery(query interfaces.Query) error {
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
    s.UnreliableNode().Send() <- packet
  }

  return nil
}

func (s *SimpleTree) OnLeave() {
}

func (s *SimpleTree) NewDHTNode() DHTNode {
  // Constructor that assumes the UnreliableNode component is filled in
  node := new(SimpleTree)

  node.Autowire(s)

  node.id       = node.UnreliableNode().Id()
  node.neighId  = node.UnreliableNode().Join()
  node.store    = make(map[string]bool)

  return node
}

func (s *SimpleTree) Key() string {
  return RandomKey()
}
