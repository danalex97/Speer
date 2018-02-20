package overlay

import (
  "testing"
  "math"
  "math/rand"
  "time"
  "github.com/danalex97/Speer/underlay"
  "github.com/danalex97/Speer/events"
)

func testUnderlayChan(nodes int) (string, string, Bridge, Bridge) {
  rand.Seed(time.Now().UTC().UnixNano())
  edges := int(math.Log2(float64(nodes)) * float64(nodes) / 2 + 1)
  network := underlay.NewRandomUniformNetwork(nodes, edges, 2, 10)

  netmap := NewNetworkMap(network)
  id1 := netmap.NewId()
  id2 := netmap.NewId()

  simulation := events.NewLazySimulation()
  netSim := underlay.NewNetworkSimulation(simulation, network)
  go netSim.Run()

  bridge1 := NewUnderlayChan(id1, netSim, netmap)
  bridge2 := NewUnderlayChan(id2, netSim, netmap)

  return id1, id2, bridge1, bridge2
}

func clearUnderChan(b Bridge) {
  b.(*UnderlayChan).simulation.Stop()
}

func TestBridgePacketDelivery(t *testing.T) {
  id1, id2, bridge1, bridge2 := testUnderlayChan(10)

  for i := 0; i < 10; i++ {
    packet := NewPacket(id1, id2, nil)
    bridge1.Send() <- packet

    recvPacket := (<-bridge2.Recv()).(Packet)
    assertEqual(t, packet.Src(), recvPacket.Src())
    assertEqual(t, packet.Dest(), recvPacket.Dest())
    assertEqual(t, packet.Payload(), recvPacket.Payload())
  }

  clearUnderChan(bridge1)
}
