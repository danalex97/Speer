package examples

import (
  . "github.com/danalex97/Speer/sdk/go"
  . "github.com/danalex97/Speer/model"
  . "github.com/danalex97/Speer/capacity"
  "runtime"
  "sync"
  "fmt"
)

type SimpleTorrent struct {
  AutowiredTorrentNode
  sync.Mutex

  id      string
  ids     []string
  links   map[string]Link
}

type controlMsg struct {
  ids []string
}

/* Interface functions. */
func (s *SimpleTorrent) OnJoin() {
  go func() {
    for {
      select {
      case _, ok := <-s.UnreliableNode().Recv():
        if ok {
          fmt.Println("Packet receive!")
        }
      case m, ok := <-s.Transfer().ControlRecv():
        if !ok {
          continue
        }
        msg := m.(controlMsg)

        s.updateIds(msg.ids)
      default:
        runtime.Gosched()
      }
    }
  }()

  go func() {

  }()
}

func (s *SimpleTorrent) OnQuery(query DHTQuery) error {
  return nil
}

func (s *SimpleTorrent) OnLeave() {
}

func (s *SimpleTorrent) NewDHTNode() DHTNode {
  // Constructor that assumes the UnreliableNode component is filled in
  node := new(SimpleTorrent)

  node.Autowire(s)

  node.id   = node.UnreliableNode().Id()
  node.ids  = []string{node.id, node.UnreliableNode().Join()}

  return node
}

func (s *SimpleTorrent) Key() string {
  return RandomKey()
}

/* Local functions */
func (s *SimpleTorrent) updateIds(ids []string) {
  allIds := make(map[string]bool)
  for _, id := range ids {
    allIds[id] = true
  }
  for _, id := range s.ids {
    allIds[id] = true
  }

  s.ids = []string{}
  for id, _ := range allIds {
    s.ids = append(s.ids, id)

    // register link if not registered
    // if _, ok := s.links[id]; !ok {
    //   s.links[id] = s.Transfer().Connect
    // }
  }
}
