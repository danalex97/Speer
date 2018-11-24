package config

import (
  "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/Speer/sdk/go"
)

func NewSimulation(template interface {}, config *Config) interfaces.ISimulation {
  builder := sdk.NewDHTSimulationBuilder(template).
    WithPoissonProcessModel(2, 2).
    WithInternetworkUnderlay(
      config.TransitDomains,
      config.TransitDomainSize,
      config.StubDomains,
      config.StubDomainSize)

  if config.Parallel {
    builder = builder.WithParallelSimulation()
  }

  capBuilder := builder.
    WithDefaultQueryGenerator().
    WithLimitedNodes(config.Nodes + 1).
    //====================================
    WithCapacities().
    WithTransferInterval(
      config.TransferInterval)

  if config.Latency {
    capBuilder = capBuilder.WithLatency()
  }

  for _, tuple := range config.CapacityNodes {
    capBuilder = capBuilder.WithCapacityNodes(
      tuple.Number,
      tuple.Upload,
      tuple.Download)
  }

  return capBuilder.Build()
}
