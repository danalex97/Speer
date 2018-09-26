package main

import (
  . "github.com/danalex97/Speer/sdk/go"
  . "github.com/danalex97/Speer/examples"
  "github.com/danalex97/Speer/interfaces"
  "github.com/danalex97/Speer/sdk/bridge"

  errLog "log"
  "runtime/pprof"
  "runtime"
  "os/signal"
  "flag"

  "math/rand"
  "time"
  "fmt"
  "os"
)

var cpuprofile = flag.String("cpuprofile", "", "Write cpu profile to `file`.")
var memprofile = flag.String("memprofile", "", "Write memory profile to `file`.")

var torrent = flag.Bool("torrent", false, "Torrent simulation.")
var metrics = flag.Bool("metrics", false, "Write metrics.")
var secs    = flag.Int("time", 10, "The time to run the simulation.")

func makeMemprofile() {
  // Profiling
  if *memprofile != "" {
    f, err := os.Create(*memprofile)
    if err != nil {
        errLog.Fatal("could not create memory profile: ", err)
    }
    runtime.GC() // get up-to-date statistics
    if err := pprof.WriteHeapProfile(f); err != nil {
        errLog.Fatal("could not write memory profile: ", err)
    }
    f.Close()
  }
}

func main() {
  env := bridge.NewEnviron("python", "python3 sdk/python/environ.py")

  nt := new(bridge.BridgedTorrent)
  ss := NewDHTSimulationBuilder(nt).
    WithPoissonProcessModel(2, 2).
    WithInternetworkUnderlay(10, 20, 20, 50).
    WithDefaultQueryGenerator().
    WithLimitedNodes(100).
    WithCapacities().
    WithLatency().
    WithTransferInterval(10).
    WithCapacityNodes(100, 10, 20).
    WithCapacityNodes(100, 30, 30).
    WithEnviron(env).
    Build()

  ss.Run()
  time.Sleep(time.Second * 1)
  ss.Stop()
  env.Stop()

  return

  rand.Seed(time.Now().UTC().UnixNano())

  // Parsing the flags
  flag.Parse()

  // Profiling
  if *cpuprofile != "" {
    f, err := os.Create(*cpuprofile)
    if err != nil {
        errLog.Fatal("could not create CPU profile: ", err)
    }
    if err := pprof.StartCPUProfile(f); err != nil {
        errLog.Fatal("could not start CPU profile: ", err)
    }
    defer pprof.StopCPUProfile()
  }
  defer makeMemprofile()
  // Get profile even on signal
  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt)
  go func(){
    for sig := range c {
      fmt.Println("Singal received:", sig)
      if *cpuprofile != "" {
        pprof.StopCPUProfile()
      }
      makeMemprofile()

      os.Exit(0)
    }
  }()

  var s interfaces.ISimulation
  if *torrent {
    nodeTemplate := new(SimpleTorrent)
    b := NewDHTSimulationBuilder(nodeTemplate).
      WithPoissonProcessModel(2, 2).
      WithInternetworkUnderlay(10, 20, 20, 50).
      WithDefaultQueryGenerator().
      WithLimitedNodes(100)
    if *metrics {
      b = b.WithMetrics()
    }
    s = b.
      WithCapacities().
      WithLatency().
      WithTransferInterval(10).
      WithCapacityNodes(100, 10, 20).
      WithCapacityNodes(100, 30, 30).
      Build()
  } else {
    nodeTemplate := new(SimpleTree)
    b := NewDHTSimulationBuilder(nodeTemplate).
      WithPoissonProcessModel(2, 2).
      WithRandomUniformUnderlay(1000, 5000, 2, 10).
      WithParallelSimulation().
      WithDefaultQueryGenerator().
      WithLimitedNodes(100)
    if *metrics {
      b = b.WithMetrics()
    }
    s = b.Build()
  }

  s.Run()

  time.Sleep(time.Second * time.Duration(*secs))
  fmt.Println("Done")
  s.Stop()

  os.Exit(0)
}
