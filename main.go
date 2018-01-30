package main

import (
  . "github.com/danalex97/Speer/sdk/go"
  . "github.com/danalex97/Speer/model"
  "github.com/danalex97/Speer/overlay"
  "runtime"
  "sync"
  "math/rand"
  "time"
  "fmt"
  "os"
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

func (s *SimpleTree) OnQuery(query DHTQuery) error {
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
    // fmt.Println(packet)
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
  const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

  b := make([]byte, 30)
  for i := range b {
    b[i] = letterBytes[rand.Int63() % int64(len(letterBytes))]
  }
  return string(b)
}

func main() {
  rand.Seed(time.Now().UTC().UnixNano())

  nodeTemplate := new(SimpleTree)
  s := NewDHTSimulationBuilder(nodeTemplate).
    WithPoissonProcessModel(2, 2).
    WithRandomUniformUnderlay(10000, 70000, 2, 10).
    WithDefaultQueryGenerator().
    WithMetrics().
    Autowire().
    Build()

  s.Run()

  time.Sleep(time.Second * 10)
  fmt.Println("Done")
  s.Stop()

  os.Exit(0)
}
