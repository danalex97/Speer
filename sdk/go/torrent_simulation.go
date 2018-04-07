package sdk

import (
  . "github.com/danalex97/Speer/capacity"
  "github.com/danalex97/Speer/events"
  "github.com/danalex97/Speer/interfaces"
)

type TorrentSimulation struct {
  *DHTSimulation

  scheduler  Scheduler
  toRegister []*registerEntry
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
    s,
    nil,
    []*registerEntry{},
  }
}

func NewTorrentSimulationBuilder(b *DHTSimulationBuilder) *TorrentSimulationBuilder {
  builder := new(TorrentSimulationBuilder)

  builder.sim = NewTorrentSimulation(b.Build())
  builder.sim.constructor = NewAutowiredTorrentNode
  builder.sim.simulation  = builder.sim

  return builder
}

func (b *TorrentSimulationBuilder) WithTransferInterval(interval int) *TorrentSimulationBuilder {
  b.sim.scheduler = NewScheduler(interval)

  // Schedule the first interval
  b.sim.underlaySimulation.Push(events.NewEvent(0, nil, b.sim.scheduler))

  return b
}

func (b *TorrentSimulationBuilder) WithCapacityNodes(nodes int, download int, upload int) *TorrentSimulationBuilder {
  b.sim.toRegister = append(b.sim.toRegister, &registerEntry{
    nodes,
    upload,
    download,
  })

  return b
}

func (b *TorrentSimulationBuilder) Build() *TorrentSimulation {
  return b.sim
}

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

  // Set connection callback
  newEngine.SetConnectCallback(func (l interfaces.Link) {
    s.scheduler.RegisterLink(l)
  })

  return newEngine
}
