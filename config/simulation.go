package config

import (
  "github.com/danalex97/Speer/config/stub"

  "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/Speer/sdk/go"

  "fmt"
  "os"
  "os/exec"
)

func NewSimulation(config *Config) interfaces.ISimulation {
  defer func() {
    if err := recover(); err != nil {
      RemoveTemplate()
      panic(err)
    }
  }()

  fmt.Println(config.Entry)
  if !TemplateExists() {
    CreateTemplate(config)

    pwd, _ := os.Getwd()
    src := fmt.Sprintf("%s/main.go", pwd)

    // run again main
    args := os.Args[1:]
    args = append(args, "run")
    args = append(args, src)
    cmd := exec.Command("go", args...)

    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    if err := cmd.Start(); err != nil {
      panic(err)
    }
    cmd.Wait()

    os.Exit(0)
  }

  defer RemoveTemplate()
  template := stub.NewNode()

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
