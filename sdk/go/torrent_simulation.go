package sdk

import (
  . "github.com/danalex97/Speer/capacity"
  "github.com/danalex97/Speer/events"
)

type TorrentSimulation struct {
  *DHTSimulation

  scheduler  Scheduler
  engines    map[DHTNode]Engine
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
    make(map[DHTNode]Engine),
    []*registerEntry{},
  }
}

func NewTorrentSimulationBuilder(b *DHTSimulationBuilder) *TorrentSimulationBuilder {
  builder := new(TorrentSimulationBuilder)
  builder.sim = NewTorrentSimulation(b.Build())
  return builder
}

func (b *TorrentSimulationBuilder) WithTransferInterval(interval int) *TorrentSimulationBuilder {
  b.sim.scheduler = NewScheduler(interval)

  // Register node update callback
  b.sim.scheduler.SetCallback(func() {
    b.sim.updateEngines()
  })

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

func (s *TorrentSimulation) updateEngines() {
  registerEngine := func () bool {
    idx      := 0
    register := s.toRegister[idx]

    register.nodes -= 1
    if register.nodes == 0 {
      s.toRegister = s.toRegister[idx + 1:]
    }

    newEngine := NewTransferEngine(register.up, register.down)
    for _, node := range s.nodeMap {
      if _, ok := s.engines[node]; !ok {
        s.engines[node] = newEngine

        // We autowire the engine
        node.(TorrentNode).autowireEngine(newEngine)

        return true
      }
    }
    return false
  }

  if len(s.toRegister) == 0 {
    return
  }

  ctr := 0
  for registerEngine() {
    ctr = ctr + 1
    if len(s.toRegister) == 0 {
      break
    }
  }
}
