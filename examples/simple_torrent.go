package examples

import (
  . "github.com/danalex97/Speer/sdk/go"
  . "github.com/danalex97/Speer/model"
  "runtime"
  "sync"
  "fmt"
)

type SimpleTorrent struct {
  AutowiredTorrentNode
  sync.Mutex

  id      string
  neighId string
}

func (s *SimpleTorrent) OnJoin() {
  go func() {
    for {
      select {
      case _, ok := <-s.UnreliableNode().Recv():
        if ok {
          fmt.Println("Packet receive!")
        }
      default:
        runtime.Gosched()
      }
    }
  }()
}

func (s *SimpleTorrent) OnQuery(query DHTQuery) error {
  return nil
}

func (s *SimpleTorrent) OnLeave() {
}

func (s *SimpleTorrent) NewDHTNode() DHTNode {
  // Constructor that assumes the UnreliableNode component is filled in
  node := new(SimpleTree)

  node.Autowire(s)

  node.id       = node.UnreliableNode().Id()
  node.neighId  = node.UnreliableNode().Join()

  return node
}

func (s *SimpleTorrent) Key() string {
  return RandomKey()
}
