package sdk

import (
  "github.com/danalex97/Speer/underlay"
  "github.com/danalex97/Speer/events"
  "github.com/danalex97/Speer/overlay"
  "github.com/danalex97/Speer/model"
)

type DHTSimulation struct {
  underlaySimulation *underlay.NetworkSimulation
  timeModel          model.TimeModel
  queryGenerator     model.DHTQueryGenerator
  node               DHTNode

  el                 eventLooper
  ql                 queryLooper
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

  return builder
}

func (b *DHTSimulationBuilder) WithPoissonProcessModel(
    arrivalRate float64,
    queryRate float64) *DHTSimulationBuilder {
  b.sim.timeModel = NewPoissonProcessModel(arrivalRate, queryRate)
  return b
}

func (b *DHTSimulationBuilder) WithDefaultQueryGenerator(
    ) *DHTSimulationBuilder {
  if b.sim.underlaySimulation == nil {
    panic("Underlay simulation component has to be appended first")
  }

  bootstrap := overlay.GetBootstrap(b.sim.underlaySimulation)
  b.sim.queryGenerator = NewDHTLedger(bootstrap)

  return b
}

func (b *DHTSimulationBuilder) WithRandomUniformUnderlay(
    nodes, edges, minLatency, maxLatency int
  ) *DHTSimulationBuilder {

  network := underlay.NewRandomUniformNetwork(nodes, edges, minLatency, maxLatency)
  s := underlay.NewNetworkSimulation(events.NewLazySimulation(), network)

  b.sim.underlaySimulation = s

  return b;
}

func (b *DHTSimulationBuilder) Build() DHTSimulation {
  if b.sim.underlaySimulation == nil {
    panic("Underlay simulation component has to be appended to build")
  }
  if b.sim.timeModel == nil {
    panic("Time model component has to be appended to build")
  }
  if b.sim.queryGenerator == nil {
    panic("Query generator component has to be appended to build")
  }

  sim := b.sim;
  b.sim = nil;

  return sim;
}

struct eventLooper {}
func (gen *eventLooper) Receive(e *Event) {
  e.payload.(Simulation).generateEvents()
}

struct queryLooper {}
func (gen *queryLooper) Receive(e *Event) {
  e.payload.(Simulation).generateQueries()
}

func (s *Simulation) generateEvents() {
  // for the moment we will only model joins
  newNode := s.node.NewDHTNode()
  newNode.OnJoin()

  // generate the next event to be handled
  event := NewEvent(
    s.underlaySimulation.Time() + int(e.timeModel.NextArrival()),
    s,
    s.eventLooper
  )
  s.Push(event)
}

func (s *Simulation) generateQueries() {
  // generate queries

  // generate the next event to be handled
  event := NewEvent(
    s.underlaySimulation.Time() + int(e.timeModel.NextQuery()),
    s,
    s.queryLooper
  )
  s.Push(event)
}

func (s *Simulation) Run() {
  s.generateEvents()
  s.generateQueries()
  go s.underlaySimulation.Run()
}
