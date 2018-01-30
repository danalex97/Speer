package main

import (
  . "github.com/danalex97/Speer/sdk/go"
  . "github.com/danalex97/Speer/model"
  "math/rand"
  "time"
)

type SimpleChain struct {
  AutowiredDHTNode

  id             string
  neigh_id       string
}

func (s *SimpleChain) OnJoin() {
}

func (s *SimpleChain) OnQuery(query DHTQuery) error {
  return nil
}

func (s *SimpleChain) OnLeave() {
}

func (s *SimpleChain) NewDHTNode() DHTNode {
  // Constructor that assumes the UnreliableNode component is filled in
  node := new(SimpleChain)

  node.Autowire(s)

  node.id       = node.UnreliableNode().Id()
  node.neigh_id = s.UnreliableNode().Id()

  return node
}

func (s *SimpleChain) Key() string {
  const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

  b := make([]byte, 30)
  for i := range b {
    b[i] = letterBytes[rand.Int63() % int64(len(letterBytes))]
  }
  return string(b)
}

func main() {
  rand.Seed(time.Now().UTC().UnixNano())

  nodeTemplate := new(SimpleChain)
  s := NewDHTSimulationBuilder(nodeTemplate).
    WithPoissonProcessModel(40, 4).
    WithRandomUniformUnderlay(10000, 70000, 2, 10).
    WithDefaultQueryGenerator().
    Autowire().
    Build()
  s.Run()

  time.Sleep(time.Second)
}
