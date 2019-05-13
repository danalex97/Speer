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
	builder := sdk.NewSimulationBuilder(template)

	if config.Network != nil {
		network := config.Network
		if network.TransitDomains == 0 || network.TransitDomainSize == 0 {
			panic("Transit domain number or transit domain size not provided or zero.")
		}
		builder = builder.WithInternetworkUnderlay(
			int(network.TransitDomains),
			int(network.TransitDomainSize),
			int(network.StubDomains),
			int(network.StubDomainSize))
	}

	if config.Parallel {
		builder = builder.WithParallelSimulation()
	}

	if config.Nodes == 0 {
		panic("Number of nodes was not provided or is 0.")
	}
	if config.LogFile != "" {
		builder = builder.WithLogs("log.json")
	}

	builder = builder.
		WithFixedNodes(int(config.Nodes))

	if config.CapacityNodes != nil && config.TransferInterval == 0 {
		panic("No transfer interval provided or transfer interval zero.")
	} else if config.CapacityNodes != nil {
		builder = builder.
			WithCapacityScheduler(int(config.TransferInterval))
	}

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

		src := fmt.Sprintf("%s/speer.go", speer)

		// run again main
		args := []string{}
		args = append(args, "run")
		args = append(args, src)
		for _, arg := range os.Args[1:] {
			args = append(args, arg)
		}
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
