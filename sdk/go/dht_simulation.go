package sdk

import (
  "github.com/danalex97/Speer/events"
  "github.com/danalex97/Speer/underlay"
  "github.com/danalex97/Speer/overlay"
  "github.com/danalex97/Speer/model"
)

type DHTSimulation struct {
  underlaySimulation *underlay.NetworkSimulation
  timeModel          model.TimeModel
  queryGenerator     model.DHTQueryGenerator
  node               DHTNode

  el                 *eventLooper
  ql                 *queryLooper
  nodeMap            map[string]DHTNode
}

type DHTSimulationBuilder struct {
  sim *DHTSimulation
}

func NewDHTSimulationBuilder(node DHTNode) *DHTSimulationBuilder {
  builder := new(DHTSimulationBuilder)
  builder.sim = new(DHTSimulation)
  builder.sim.node = node

  builder.sim.el   = new(eventLooper)
  builder.sim.ql   = new(queryLooper)
  builder.sim.nodeMap = make(map[string]DHTNode)

  return builder
}

func (b *DHTSimulationBuilder) WithPoissonProcessModel(
    arrivalRate float64,
    queryRate float64) *DHTSimulationBuilder {
  b.sim.timeModel = model.NewPoissonProcessModel(arrivalRate, queryRate)
  return b
}

func (b *DHTSimulationBuilder) WithDefaultQueryGenerator(
    ) *DHTSimulationBuilder {
  if b.sim.underlaySimulation == nil {
    panic("Underlay simulation component has to be appended first")
  }

  bootstrap := overlay.GetBootstrap(b.sim.underlaySimulation)
  b.sim.queryGenerator = model.NewDHTLedger(bootstrap)

  return b
}

func (b *DHTSimulationBuilder) WithRandomUniformUnderlay(
    nodes int,
    edges int,
    minLatency int,
    maxLatency int) *DHTSimulationBuilder {
  network := underlay.NewRandomUniformNetwork(nodes, edges, minLatency, maxLatency)
  s := underlay.NewNetworkSimulation(events.NewLazySimulation(), network)

  b.sim.underlaySimulation = s

  return b;
}

func (b *DHTSimulationBuilder) Autowire() *DHTSimulationBuilder{
  aw := b.sim.node.autowire().(*AutowiredDHTNode)
  aw.node = overlay.NewUnreliableSimulatedNode(b.sim.underlaySimulation)
  return b
}

func (b *DHTSimulationBuilder) Build() *DHTSimulation {
  if b.sim.underlaySimulation == nil {
    panic("Underlay simulation component has to be appended to build")
  }
  if b.sim.timeModel == nil {
    panic("Time model component has to be appended to build")
  }
  if b.sim.queryGenerator == nil {
    panic("Query generator component has to be appended to build")
  }
  if b.sim.node == nil {
    panic("Node protocol component has to be appended to build")
  }

  sim := b.sim;
  b.sim = nil;

  return sim;
}

type eventLooper struct {}
func (gen *eventLooper) Receive(e *events.Event) *events.Event {
  e.Payload().(*DHTSimulation).generateEvents()
  return nil
}

type queryLooper struct {}
func (gen *queryLooper) Receive(e *events.Event) *events.Event {
  e.Payload().(*DHTSimulation).generateQueries()
  return nil
}

func (s *DHTSimulation) generateEvents() {
  // for the moment we will only model joins
  newNode := s.node.NewDHTNode()
  // id selection should probabily be moved to SDK (?)
  // now the overlay sits somewhere between the transport and netowrk layer
  id      := newNode.UnreliableNode().Id()
  s.nodeMap[id] = newNode
  newNode.OnJoin()

  // generate the next event to be handled
  time := s.underlaySimulation.Time() + int(s.timeModel.NextArrival())
  event := events.NewEvent(time, s, s.el)
  s.underlaySimulation.Push(event)
}

func (s *DHTSimulation) generateQueries() {
  // generate queries
  query := s.queryGenerator.Next()
  // deliver queries to nodes as well

  // the template node is not in the map, so we need to avoid it if possible
  // TODO: need to fix this bug, as the bootstrap may break!
  if node, ok := s.nodeMap[query.Node()]; ok {
    go node.OnQuery(*query)
  }

  // generate the next event to be handled
  time := s.underlaySimulation.Time() + int(s.timeModel.NextQuery())
  event := events.NewEvent(time, s, s.ql)

  s.underlaySimulation.Push(event)
}

func (s *DHTSimulation) Run() {
  s.generateEvents()
  s.generateQueries()
  go s.underlaySimulation.Run()
}

func (s *DHTSimulation) Stop() {
  s.underlaySimulation.Stop()
}
