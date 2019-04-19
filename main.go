package main

import (
	"github.com/danalex97/Speer/config"

	"flag"
	errLog "log"
	"os/signal"
	"runtime"
	"runtime/pprof"

	"fmt"
	"math/rand"
	"os"
	"time"
)

var cpuprofile = flag.String("cpuprofile", "", "Write cpu profile to `file`.")
var memprofile = flag.String("memprofile", "", "Write memory profile to `file`.")
var configPath = flag.String("config", "./examples/config/broadcast.json", "Path to configuration file.")

var secs = flag.Int("time", 1, "The time to run the simulation.")

func makeCPUProfile() func() {
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			errLog.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			errLog.Fatal("could not start CPU profile: ", err)
		}
		return pprof.StopCPUProfile
	}
	return func() {}
}

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

func setSignals() {
	// Get profile even on signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Println("Signal received:", sig)
			if *cpuprofile != "" {
				pprof.StopCPUProfile()
			}
			makeMemprofile()

			os.Exit(0)
		}
	}()
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	// Parsing the flags
	flag.Parse()

	// Profiling
	defer makeCPUProfile()()
	defer makeMemprofile()
	setSignals()

	jsonConfig := config.JSONConfig(*configPath)
	simulation := config.NewSimulation(jsonConfig)
	simulation.Run()

	time.Sleep(time.Second * time.Duration(*secs))
	simulation.Stop()
}
