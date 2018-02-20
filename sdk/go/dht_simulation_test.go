package sdk

import (
  . "github.com/danalex97/Speer/model"

  "testing"
  "math/rand"
  "sync"
  "time"
  "fmt"
  "runtime"
)

var joins   int = 0
var queries int = 0

type mockNode struct {
  AutowiredDHTNode
  sync.Mutex

  id           string
  neighId      string
}

func (s *mockNode) OnJoin() {
  s.Lock()
  defer s.Unlock()

  joins += 1
  runtime.Gosched()
}

func (s *mockNode) OnQuery(query DHTQuery) error {
  s.Lock()
  defer s.Unlock()

  queries += 1
  return nil
}

func (s *mockNode) OnLeave() {
}

func (s *mockNode) NewDHTNode() DHTNode {
  // Constructor that assumes the UnreliableNode component is filled in
  node := new(mockNode)

  node.Autowire(s)

  node.id       = node.UnreliableNode().Id()
  node.neighId  = node.UnreliableNode().Join()

  return node
}

func (s *mockNode) Key() string {
  return ""
}

func TestGoSDKQueriesGetGenerated(t *testing.T) {
  // WARN: The query rates offer no guarantees due to packet exanges
  // WARN: System doesn't progress on size limit obtained.
  rand.Seed(time.Now().UTC().UnixNano())

  arrRate := 40.0
  // WARN: Small query rates yield big variance, considering virtual time
  // in terms of millisconds would help
  queryRate := 10.0

  nodeTemplate := new(mockNode)
  s := NewDHTSimulationBuilder(nodeTemplate).
    WithPoissonProcessModel(arrRate, queryRate).
    WithRandomUniformUnderlay(1000, 70000, 2, 10).
    WithDefaultQueryGenerator().
    Autowire().
    Build()

  go s.Run()

  time.Sleep(time.Millisecond * 1500)
  s.Stop()

  virtualTime := s.Time()

  fmt.Println("Go SDK test: (joins: %s, queries: %s, virtualTime: %s",
    joins, queries, virtualTime)
  if float64(joins) < 0.9 * float64(virtualTime) / arrRate {
    t.Fatalf("Unexpectedly small number of joins: %s; expected at %s",
      joins, 0.9 * float64(virtualTime) * arrRate)
  }
  if float64(queries) < 0.9 * float64(virtualTime) / queryRate {
    t.Fatalf("Unexpectedly small number of queries: %s; expected at %s",
      queries,  0.9 * float64(virtualTime) * queryRate)
  }
  if float64(joins) > 1.1 * float64(virtualTime) / arrRate {
    t.Fatalf("Unexpectedly big number of joins: %s; expected at %s",
      joins, 1.1 * float64(virtualTime) * arrRate)
  }
  if float64(queries) > 1.1 * float64(virtualTime) / queryRate {
    t.Fatalf("Unexpectedly big number of queries: %s; expected at %s",
      queries,  1.1 * float64(virtualTime) * queryRate)
  }
}
