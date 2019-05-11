package overlay

import (
	"github.com/danalex97/Speer/events"
	"github.com/danalex97/Speer/underlay"

	"math"
	"math/rand"
	"testing"
	"time"
)

func testUnderlayChan(nodes int) (string, string, LatencyConnector, LatencyConnector) {
	rand.Seed(time.Now().UTC().UnixNano())
	edges := int(math.Log2(float64(nodes))*float64(nodes)/2 + 1)
	network := underlay.NewRandomUniformNetwork(nodes, edges, 2, 10)

	netmap := NewNetworkMap(network)
	id1 := netmap.NewId()
	id2 := netmap.NewId()

	simulation := events.NewLazySimulation()
	netSim := underlay.NewNetworkSimulation(simulation, network)
	go netSim.Run()

	time.Sleep(100 * time.Millisecond)

	bridge1 := NewUnderlayChan(id1, netSim, netmap)
	bridge2 := NewUnderlayChan(id2, netSim, netmap)

	return id1, id2, bridge1, bridge2
}

func clearUnderChan(b LatencyConnector) {
	b.(*UnderlayChan).simulation.Stop()
}

func TestLatencyConnectorPacketDelivery(t *testing.T) {
	_, id2, bridge1, bridge2 := testUnderlayChan(10)

	for i := 0; i < 10; i++ {
		t.Logf("LatencyConnector packet delivery test -- packet %d\n", i)
		bridge1.ControlSend(id2, "message")
		assertEqual(t, "message", <-bridge2.ControlRecv())
	}

	clearUnderChan(bridge1)
}

func TestLatencyConnectorSendPacketToSelf(t *testing.T) {
	done := make(chan bool, 1)

	go func() {
		id1, _, bridge1, _ := testUnderlayChan(10)

		bridge1.ControlSend(id1, "message")
		assertEqual(t, "message", <-bridge1.ControlRecv())

		done <- true
	}()

	time.Sleep(200 * time.Millisecond)
	select {
	case <-done:
	default:
		t.Fatalf("Test timeout.")
	}
}
