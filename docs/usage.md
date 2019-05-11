### SDKs

- Golang
- Python [[planned](roadmap.md)]

### Golang SDK

The current implementation provides a Golang SDK. In order to be able to use the simulator, the user can implement the **TorrentNode** interface. (DHTNode interface is deperecated)

Each interface comes with a **Util** interface provided to the user for interacting with a simulation:
```go
type TorrentNode interface {
    // The general method that is just a runner.
    OnJoin()

    // A method that should be called when a node leaves the network.
    OnLeave()
}

type TorrentNodeConstructor interface {
    New(util TorrentNodeUtil) TorrentNode
}
```

A user needs to implement all the methods provided below to obtain a valid node to use in a simulation:
```go
import (
  . "github.com/danalex97/Speer/interfaces"
)

type ExampleTorrent struct {
  //[...]
}

func (s *ExampleTorrent) OnJoin() {
  //[...]
}

func (s *ExampleTorrent) OnLeave() {
  //[...]
}

func (s *ExampleTorrent) New(util TorrentNodeUtil) TorrentNode {
  //[...]
}
```

The `New` function is used to generate new nodes from an empty structure template. The structure `TorrentNodeUtil` provides a set of functions which `ExampleTorrent` can use to interact with the simulation:
```go
type TorrentNodeUtil interface {
    Transport() Transport // Interface used to send data to other nodes.

    Id()   string  // The ID of the node.
    Join() string  // The ID of another node in the network.

    Time() func() int // Function that can be used to retrieve the simulation global virtual time.
}
```

Some examples on how to use the **TorrentNodeUtil** can be found in the **Examples** section below.

### Running simulation

#### Using the Speer configuration stubs

Speer allows running a simulation directly from a JSON configuration that provides simulation parameters and a **[path-to-implementation-package]/[name-of-structure-to-simulate]**. The simplest way to run a simulation is to run:
```
go run speer.go -config=[configuration-path]
```

The configuration has to provide the path to your source code that implements the `TorrentNode` and `TorrentNodeConstructor` interfaces. For an example, check `examples/config/big.json`.

#### Using a JSON configuration

```go
package examples

import (
  . "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/Speer/config"
)

[...]

func example() {
	jsonConfig := config.JSONConfig("./config/big_no_entry.json")

    // Example has to implement the `TorrentNodeUtil` interface
	template := Example.New()

	simulation := NewSimulationFromTemplate(config, template)

	simulation.Run()
	time.Sleep(time.Second * time.Duration(10))
	simulation.Stop()
}
```

For an example of a JSON configuration check `examples/config/big_no_entry.json`.

#### Using `sdk/go`

To SDK packet can be used to build and run a simulation. The SDK packet offers a builder interface which allows building custom simulations. An example is:
```go
import (
  . "github.com/danalex97/Speer/sdk/go"
)

// [...]

func main() {
  // We need to provide a node template which implements the
  // TorrentNode interface.
  nodeTemplate := new(ExampleTorrent)

  // The simulation is created via a builder pattern.
  // If an options is missing, the builder will panic.
  s := NewDHTSimulationBuilder(nodeTemplate).
    WithPoissonProcessModel(2, 2).
    WithInternetworkUnderlay(10, 20, 20, 50).
    WithDefaultQueryGenerator().
    WithLimitedNodes(100).
    // The WithCapacities options are allowed when we use
    // TorrentNodes instead of DHTNodes.
    WithCapacities().
    WithLatency().
    WithTransferInterval(10).
    WithCapacityNodes(100, 10, 20).
    WithCapacityNodes(100, 30, 30).
    Build()

  // Starts the simulation in a different Goroutine.
  s.Run()

  //[...]

  // Sends a stop signal to the simulation. It will stop the simulation gracefully as soon as possible.
  s.Stop()
}
```

To see a full list of the options you can consult the following links:
- [DHTSimulation](https://godoc.org/github.com/danalex97/Speer/sdk/go#DHTSimulationBuilder)
- [TorrentSimulation](https://godoc.org/github.com/danalex97/Speer/sdk/go#TorrentSimulationBuilder)

A **DHTSimulationBuilder** can be converted into a **TorrentSimulationBuilder** by using the method **WithCapacities()** as shown in the example above.


### Examples

We will show some examples on how to use the **TorrentNodeUtil**:
- sending a message:

```go
util.Transport().ControlSend(util.Join(), "message")
```
- receiving a message:

```go
msg := <-util.Transport().ControlRecv()
```
- pinging a node:

```go
if util.Transport().ControlPing(util.Join()) {
    //[...]
}
```
- sending a big amount of data via a link:

```go
// Creating the link
link := util.Transport().Connect(util.Join())

// Sending the data
link.Upload(Data{
  Id   : "someUniqueId", // Some ID associated with the message.
                         // The ID can be used for keeping the actual data or metadata.
  Size : 1000, // Total data size in bits.
})
```

- getting data from a link:

```go
data := <-link.Download()
```

For more examples on how to write the code for a node, check the **examples folder**. For a more complex example, check this [repository](https://github.com/danalex97/nfsTorrent).

### Running examples

The see the options provided by `speer.go` run:
```
go run speer.go -h
```

To run the *TorrentNode* example from the examples folder:
```
go run speer.go -config=./examples/config/big.json
```

The *DHTNode* interface will be deprecated, but one can run it by checking out to version 1.0 and running:
```
go run speer.go
```
