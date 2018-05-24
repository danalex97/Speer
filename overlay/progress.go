package overlay

import (
  "github.com/danalex97/Speer/events"
  "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/Speer/underlay"

  "sync"
)

type GroupProgress interfaces.GroupProgress

// There is one progress property per network simulation.
var progressSet = make(map[*underlay.NetworkSimulation]*TransmissionProgress)
var mutex       = new(sync.Mutex)

// The delay for pushing a packet into the network is at most 100ms.
const progressInterval int = 100

type TransmissionProgress struct {
  simulation *underlay.NetworkSimulation

  // Progress properties that regulate the transfer of packets.
  PushProgress GroupProgress
  PullProgress GroupProgress
}

func GetTransmissionProgress(simulation *underlay.NetworkSimulation) (progress *TransmissionProgress) {
  mutex.Lock()
  defer mutex.Unlock()

  if oldProgress, ok := progressSet[simulation]; ok {
    progress = oldProgress
  } else {
    progress = &TransmissionProgress{
      PushProgress : events.NewWGProgress(progressInterval),
      PullProgress : events.NewWGProgress(progressInterval),
      simulation   : simulation,
    }

    // Initialize and update the progress map.
    progress.Init()
    progressSet[simulation] = progress
  }
  return
}

func (p *TransmissionProgress) Init() {
  p.simulation.RegisterProgress(events.NewProgressProperty(
    p.PushProgress,
    progressInterval,
  ))
  p.simulation.RegisterProgress(events.NewProgressProperty(
    p.PullProgress,
    progressInterval,
  ))
}
