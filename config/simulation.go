package config

import (
	"github.com/danalex97/Speer/config/stub"

	"github.com/danalex97/Speer/interfaces"
	"github.com/danalex97/Speer/sdk/go"

	"fmt"
	"os"
	"os/exec"
)

func NewSimulationFromTemplate(
	config *Config,
	template interfaces.Node,
) interfaces.ISimulation {
	if config.TransitDomains == 0 || config.TransitDomainSize == 0 {
		panic("Transit domain number or transit domain size not provided or zero.")
	}

	builder := sdk.NewSimulationBuilder(template).
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
	if config.LogFile != "" {
		builder = builder.WithLogs("log.json")
	}

	builder = builder.
		WithFixedNodes(int(config.Nodes)).
		WithCapacityScheduler(int(config.TransferInterval))

	// [TODO] allow running without latency
	// if config.Latency {
	// 	builder = builder.WithLatency()
	// }

	for _, tuple := range config.CapacityNodes {
		builder = builder.WithCapacityNodes(
			int(tuple.Number),
			int(tuple.Upload),
			int(tuple.Download))
	}

	return builder.Build()
}

func NewSimulation(config *Config) interfaces.ISimulation {
	defer func() {
		if err := recover(); err != nil {
			RemoveTemplate()
			panic(err)
		}
	}()

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

	fmt.Println("Template:", config.Entry)

	defer RemoveTemplate()
	return NewSimulationFromTemplate(config, stub.NewNode())
}
