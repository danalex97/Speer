package main

import (
  . "github.com/danalex97/Speer/sdk/go"
  . "github.com/danalex97/Speer/examples"
  "math/rand"
  "time"
  "fmt"
  "os"
)

func main() {
  rand.Seed(time.Now().UTC().UnixNano())

  nodeTemplate := new(SimpleTorrent)
  // nodeTemplate := new(SimpleTree)
  s := NewDHTSimulationBuilder(nodeTemplate).
    WithPoissonProcessModel(2, 2).
    // WithRandomUniformUnderlay(1000, 5000, 2, 10).
    WithInternetworkUnderlay(10, 50, 20, 50).
    // WithInternetworkUnderlay(10, 50, 100, 100).
    WithDefaultQueryGenerator().
    WithLimitedNodes(100).
    // WithMetrics().
    //====================================
    WithCapacities().
    WithTransferInterval(10).
    WithCapacityNodes(100, 10, 20).
    WithCapacityNodes(100, 30, 30).
    Build()

  s.Run()

  time.Sleep(time.Second * 10)
  fmt.Println("Done")
  s.Stop()

  os.Exit(0)
}
