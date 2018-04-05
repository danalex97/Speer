package sdk

import (
  // "github.com/danalex97/Speer/capacity"
)

type TorrentSimulation struct {
  *DHTSimulation
}

type TorrentSimulationBuilder struct {
  sim *TorrentSimulation
}

func NewTorrentSimulation(s *DHTSimulation) *TorrentSimulation {
  return &TorrentSimulation{s}
}

func NewTorrentSimulationBuilder(b *DHTSimulationBuilder) *TorrentSimulationBuilder {
  builder := new(TorrentSimulationBuilder)
  builder.sim = NewTorrentSimulation(b.Build())
  return builder
}

func (b *TorrentSimulationBuilder) Build() *TorrentSimulation {
  return b.sim
}
