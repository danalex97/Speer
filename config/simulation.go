package config

import (
  "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/Speer/sdk/go"
)

func NewSimulation(config *Config) interfaces.ISimulation {
  var template interface {}
  if config.Lang == "go" {
    if config.Entry == "" {
      panic("Entry point not provided.")
    }
    template = NewGoTemplate(config.Entry)
  } else {
    panic("Lanaguage " + config.Lang + " not supported")
  }

  if config.TransitDomains == 0 || config.TransitDomainSize == 0 {
    panic("Transit domain number or transit domain size not provided or zero.")
  }

  builder := sdk.NewDHTSimulationBuilder(template).
    WithPoissonProcessModel(2, 2).
    WithInternetworkUnderlay(
      int(config.TransitDomains),
      int(config.TransitDomainSize),
      int(config.StubDomains),
      int(config.StubDomainSize))

  if config.Parallel {
    builder = builder.WithParallelSimulation()
  }

  if config.TransferInterval == 0 {
    panic("No transfer interval provided or transfer interval zero.")
  }
  if config.Nodes == 0 {
    panic("Number of nodes was not provided or is 0.")
  }

  capBuilder := builder.
    WithDefaultQueryGenerator().
    WithLimitedNodes(int(config.Nodes) + 1).
    //====================================
    WithCapacities().
    WithTransferInterval(
      int(config.TransferInterval))

  if config.Latency {
    capBuilder = capBuilder.WithLatency()
  }

  for _, tuple := range config.CapacityNodes {
    capBuilder = capBuilder.WithCapacityNodes(
      int(tuple.Number),
      int(tuple.Upload),
      int(tuple.Download))
  }

  return capBuilder.Build()
}
