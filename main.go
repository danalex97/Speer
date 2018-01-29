package main

import (
  . "github.com/danalex97/Speer/sdk/go"
  . "github.com/danalex97/Speer/model"
  "math/rand"
  "time"
)

type SimpleTree struct {
  AutowiredDHTNode

  id             string
  neigh_id       string
}

func (s *SimpleTree) OnJoin() {
}

func (s *SimpleTree) OnQuery(query DHTQuery) error {
  return nil;
}

func (s *SimpleTree) OnLeave() {
}

func (s *SimpleTree) NewDHTNode() DHTNode {
  // Constructor that assumes the UnreliableNode component is filled in
  return nil
}

func (s *SimpleTree) Key() string {
  return "key"
}

func main() {
  rand.Seed(time.Now().UTC().UnixNano())

  nodeTemplate := new(SimpleTree)
  s := NewDHTSimulationBuilder(nodeTemplate).
    WithPoissonProcessModel(60, 4).
    WithRandomUniformUnderlay(10000, 70000, 2, 10).
    WithDefaultQueryGenerator().
    Autowire().
    Build()
  s.Run()
}
