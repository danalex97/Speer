package main

import (
  "github.com/danalex97/Speer/config"

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
var configPath = flag.String("config", "./examples/config/big.json", "Configuration file.")

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

  jsonConfig := config.JSONConfig(*configPath)
  simulation := config.NewSimulation(jsonConfig)
  simulation.Run()

  time.Sleep(time.Second * time.Duration(*secs))
  fmt.Println("Done")
  simulation.Stop()

  os.Exit(0)
}
