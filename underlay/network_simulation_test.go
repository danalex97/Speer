package underlay

import (
	"math/rand"
	"testing"
	"time"
	// "fmt"
	. "github.com/danalex97/Speer/events"
)

func TestNetworkSimulationSmallNetwork(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	network := NewRandomUniformNetwork(40, 150, 1, 2)

	underSim := NewNetworkSimulation(NewLazySimulation(), network)

	for i := 0; i < LazyQueueChanSize/2; i++ {
		r1 := underSim.Network().RandomRouter()
		r2 := underSim.Network().RandomRouter()
		underSim.SendPacket(NewPacket(r1, r2, nil))
	}

	go underSim.Run()

	time.Sleep(50 * time.Millisecond)
	underSim.Stop()

	if underSim.Time() < 2 {
		t.Fatalf("Simulation too slow.")
	}
}

func BenchmarkNetworkSimulationSmallNetworkSpeed(t *testing.B) {
	rand.Seed(time.Now().UTC().UnixNano())
	network := NewRandomUniformNetwork(1000, 10000, 1, 2)

	underSim := NewNetworkSimulation(NewLazySimulation(), network)

	for i := 0; i < LazyQueueChanSize; i++ {
		r1 := underSim.Network().RandomRouter()
		r2 := underSim.Network().RandomRouter()
		underSim.SendPacket(NewPacket(r1, r2, nil))
	}

	go underSim.Run()

	time.Sleep(1000 * time.Millisecond)
	underSim.Stop()

	if underSim.Time() < 3 {
		t.Fatalf("Simulation too slow.")
	}
}

func BenchmarkNetworkSimulationBigNetworkSpeed(t *testing.B) {
	rand.Seed(time.Now().UTC().UnixNano())
	network := NewRandomUniformNetwork(10000, 100000, 1, 2)

	underSim := NewNetworkSimulation(NewLazySimulation(), network)

	for i := 0; i < LazyQueueChanSize/2; i++ {
		r1 := underSim.Network().RandomRouter()
		r2 := underSim.Network().RandomRouter()
		underSim.SendPacket(NewPacket(r1, r2, nil))
	}

	go underSim.Run()

	time.Sleep(7000 * time.Millisecond)
	underSim.Stop()

	if underSim.Time() < 4 {
		t.Fatalf("Simulation too slow.")
	}
}
