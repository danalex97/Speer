# Running a simulation

#### Using the Speer configuration stubs

Speer allows running a simulation directly from a JSON configuration that provides simulation parameters and a **[path-to-implementation-package]/[name-of-structure-to-simulate]**. The simplest way to run a simulation is to run:
```
speer -config=[configuration-path]
```

The configuration has to provide the path to your source code that implements the `Node` interface. For an example, check `examples/config/broadcast.json`.

#### Using a JSON configuration

```go
package examples

import (
  . "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/Speer/config"
)

[...]

func example() {
	jsonConfig := config.JSONConfig("./config/no_entry.json")

	// Example has to implement the `TorrentNodeUtil` interface
	template := Example.New()

	simulation := NewSimulationFromTemplate(config, template)

	simulation.Run()
	time.Sleep(time.Second * time.Duration(10))
	simulation.Stop()
}
```

For an example of a JSON configuration check `examples/config/no_entry.json`.

#### Using `sdk/go`

To SDK packet can be used to build and run a simulation. The SDK packet offers a builder interface which allows building custom simulations. An example is:
```go
import (
  . "github.com/danalex97/Speer/sdk/go"
)

// [...]

func main() {
	// We need to provide a node template which implements the
	// Node interface.
	template := new(Example)

	// The simulation is created via a builder pattern.
	// If an options is missing, the builder will panic.
	sim := NewSimulationBuilder(template).
		WithInternetworkUnderlay(5, 5, 5, 5).
		WithParallelSimulation().
		WithFixedNodes(10).
		WithCapacityScheduler(1).
		WithCapacityNodes(10, 1, 1).
		Build()

	// Starts the simulation in a different Goroutine.
	go s.Run()

	// Sends a stop signal to the simulation. It will stop the simulation gracefully after 0.1 seconds
	time.Sleep(100 * time.Millisecond)
	s.Stop()
}
```

To see a full list of the options you can consult the following [link](https://godoc.org/github.com/danalex97/Speer/sdk/go#SimulationBuilder).
