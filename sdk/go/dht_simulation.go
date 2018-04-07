package sdk

import (
  "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/Speer/events"
  "github.com/danalex97/Speer/underlay"
  "github.com/danalex97/Speer/overlay"
  "github.com/danalex97/Speer/model"
  "github.com/danalex97/Speer/metrics"
  "time"
  "fmt"
)

type DHTSimulation struct {
  underlaySimulation *underlay.NetworkSimulation
  timeModel          model.TimeModel
  queryGenerator     model.DHTQueryGenerator
  template           interface {}
  constructor        func(interfaces.UnreliableNode, interface {}) DHTNode

  el                 *eventLooper
  ql                 *queryLooper
  nodeMap            map[string]DHTNode
  nodeLimit          int
  nodes              int
}

type DHTSimulationBuilder struct {
  sim *DHTSimulation
}

const maxNodeLimit int = 10000000

func NewDHTSimulationBuilder(template interface {}) *DHTSimulationBuilder {
  builder := new(DHTSimulationBuilder)
  builder.sim = new(DHTSimulation)
  builder.sim.template = template

  builder.sim.el   = new(eventLooper)
  builder.sim.ql   = new(queryLooper)
  builder.sim.nodeMap = make(map[string]DHTNode)
  builder.sim.constructor = NewAutowiredDHTNode

  return builder
}

func (b *DHTSimulationBuilder) WithMetrics() *DHTSimulationBuilder {
  globalObserver := events.NewGlobalEventObserver()
  b.sim.underlaySimulation.RegisterObserver(globalObserver)

  netMap := overlay.GetBootstrap(b.sim.underlaySimulation).(*overlay.NetworkMap)
  metrics := metrics.NewMetrics(globalObserver, netMap)

  go metrics.Run()

  return b
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

  b.sim.nodeLimit = maxNodeLimit
  b.sim.nodes = 0

  return b
}

func (b *DHTSimulationBuilder) WithLimitedNodes(
  nodeLimit int) *DHTSimulationBuilder {
  if b.sim.queryGenerator == nil {
    panic("Query generator component has to be appended first.")
  }

  b.sim.nodeLimit = nodeLimit

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

func (b *DHTSimulationBuilder) WithInternetworkUnderlay(
    transitDomains int,
    transitDomainSize int,
    stubDomains int,
    stubDomainSize int) *DHTSimulationBuilder {
  network := underlay.NewInternetwork(transitDomains, transitDomainSize, stubDomains, stubDomainSize)
  s := underlay.NewNetworkSimulation(events.NewLazySimulation(), network)

  fmt.Printf("Internetwork built with %d nodes.\n", len(network.Routers))
  b.sim.underlaySimulation = s

  return b;
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
  if b.sim.template == nil {
    panic("Template component has to be appended to build")
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
  // this should be modified when we model leaves
  if s.nodes > s.nodeLimit {
    // this stops the event generation loop
    fmt.Println("Node generation was stopped.")
    return
  }

  // for the moment we will only model joins
  overlayNode := overlay.NewUnreliableSimulatedNode(s.underlaySimulation)
  newNode := s.constructor(overlayNode, s.template)

  // id selection should probabily be moved to SDK (?)
  // now the overlay sits somewhere between the transport and netowrk layer
  id      := newNode.UnreliableNode().Id()
  s.nodeMap[id] = newNode
  newNode.OnJoin()

  // update node count
  s.nodes += 1

  // generate the next event to be handled
  time := s.underlaySimulation.Time() + int(s.timeModel.NextArrival())
  event := events.NewEvent(time, s, s.el)

  // the log event is used only by the metrics module
  logEvent := events.NewEvent(time, *model.NewJoin(id), nil)

  s.underlaySimulation.Push(event)
  s.underlaySimulation.Push(logEvent)
}

func (s *DHTSimulation) generateQueries() {
  // generate queries
  query := s.queryGenerator.Next()
  // deliver queries to nodes as well

  // the template node is not in the map, so we need to avoid it if possible
  // TODO: need to fix this bug, as the bootstrap may break!
  if node, ok := s.nodeMap[query.Node()]; ok {
    go node.OnQuery(query)
  }

  // generate the next event to be handled
  time := s.underlaySimulation.Time() + int(s.timeModel.NextQuery())
  event := events.NewEvent(time, s, s.ql)

  // the log event is used only by the metrics module
  logEvent := events.NewEvent(time, query, nil)

  s.underlaySimulation.Push(event)
  s.underlaySimulation.Push(logEvent)
}

func (s *DHTSimulation) Time() int {
  return s.underlaySimulation.Time()
}

func (s *DHTSimulation) Run() {
  time.Sleep(time.Second * 1)

  s.generateEvents()
  s.generateQueries()
  go s.underlaySimulation.Run()
}

func (s *DHTSimulation) Stop() {
  s.underlaySimulation.Stop()
}

// Entry point for torrent simulation
func (t *DHTSimulationBuilder) WithCapacities() *TorrentSimulationBuilder {
  return NewTorrentSimulationBuilder(t)
}
