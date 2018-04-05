package sdk

import (
  . "github.com/danalex97/Speer/capacity"
  "github.com/danalex97/Speer/events"
)

type TorrentSimulation struct {
  *DHTSimulation

  scheduler Scheduler
}

type TorrentSimulationBuilder struct {
  sim *TorrentSimulation
}

func NewTorrentSimulation(s *DHTSimulation) *TorrentSimulation {
  return &TorrentSimulation{s, nil}
}

func NewTorrentSimulationBuilder(b *DHTSimulationBuilder) *TorrentSimulationBuilder {
  builder := new(TorrentSimulationBuilder)
  builder.sim = NewTorrentSimulation(b.Build())
  return builder
}

func (b *TorrentSimulationBuilder) WithTransferInterval(interval int) *TorrentSimulationBuilder {
  b.sim.scheduler = NewScheduler(interval)

  // Schedule the first interval
  b.sim.underlaySimulation.Push(events.NewEvent(0, nil, b.sim.scheduler))

  return b
}

func (b *TorrentSimulationBuilder) Build() *TorrentSimulation {
  return b.sim
}
