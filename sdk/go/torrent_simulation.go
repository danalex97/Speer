package sdk

import (
  . "github.com/danalex97/Speer/capacity"
  "github.com/danalex97/Speer/events"
  "github.com/danalex97/Speer/overlay"
  "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/Speer/sdk/bridge"
)

type TorrentSimulation struct {
  *DHTSimulation

  scheduler  Scheduler
  toRegister []*registerEntry

  // TODO: experimental
  env *bridge.Environ

  latency    bool
}

type TorrentSimulationBuilder struct {
  sim *TorrentSimulation
}

type registerEntry struct {
  nodes int
  up    int
  down  int
}

func NewTorrentSimulation(s *DHTSimulation) *TorrentSimulation {
  return &TorrentSimulation{
    DHTSimulation : s,
    scheduler     :  nil,
    toRegister    : []*registerEntry{},
    latency       : false,
  }
}

func NewTorrentSimulationBuilder(b *DHTSimulationBuilder) *TorrentSimulationBuilder {
  builder := new(TorrentSimulationBuilder)

  builder.sim = NewTorrentSimulation(b.Build())
  builder.sim.constructor = NewAutowiredTorrentNode
  builder.sim.simulation  = builder.sim

  return builder
}

func (b *TorrentSimulationBuilder) WithLatency() *TorrentSimulationBuilder {
  b.sim.latency = true
  return b
}

func (b *TorrentSimulationBuilder) WithTransferInterval(interval int) *TorrentSimulationBuilder {
  b.sim.scheduler = NewScheduler(interval)

  // Schedule the first interval
  b.sim.underlaySimulation.Push(events.NewEvent(0, nil, b.sim.scheduler))

  return b
}

// TODO: experimental
func (b *TorrentSimulationBuilder) WithEnviron(env *bridge.Environ) *TorrentSimulationBuilder {
  b.sim.env = env
  go env.Start()

  return b
}

func (b *TorrentSimulationBuilder) WithCapacityNodes(nodes int, upload int, download int) *TorrentSimulationBuilder {
  b.sim.toRegister = append(b.sim.toRegister, &registerEntry{
    nodes : nodes,
    up    : upload,
    down  : download,
  })

  return b
}

func (b *TorrentSimulationBuilder) Build() *TorrentSimulation {
  return b.sim
}

// updateEngine is called when the TorrentNode is initialized.
func (s *TorrentSimulation) updateEngine(node interfaces.UnreliableNode) Engine {
  if len(s.toRegister) == 0 {
    return nil
  }

  idx      := 0
  register := s.toRegister[idx]

  register.nodes -= 1
  if register.nodes == 0 {
    s.toRegister = s.toRegister[idx + 1:]
  }

  newEngine := NewTransferEngine(
    register.up,
    register.down,
    node.Id(),
  )
  if s.latency {
    newEngine = NewTransferLatencyEngine(
      newEngine.(*TransferEngine),
      node.(overlay.UnreliableNode),
    )
  }

  // Set connection callback
  newEngine.SetConnectCallback(func (l interfaces.Link) {
    s.scheduler.RegisterLink(l)
  })

  return newEngine
}
