package overlay

import (
  "testing"
  "math"
  "math/rand"
  "time"
  "github.com/danalex97/Speer/underlay"
  "github.com/danalex97/Speer/events"
)

func testUnderlayNetsim(nodes int) *underlay.NetworkSimulation {
  rand.Seed(time.Now().UTC().UnixNano())
  edges := int(math.Log2(float64(nodes)) * float64(nodes) / 2 + 1)
  network := underlay.NewRandomUniformNetwork(nodes, edges, 2, 10)

  simulation := events.NewLazySimulation()
  netSim := underlay.NewNetworkSimulation(simulation, network)

  return netSim
}

func TestUnreliableNodesPacketSending(t *testing.T) {
  netsim := testUnderlayNetsim(10)

  n1 := NewUnreliableSimulatedNode(netsim)
  n2 := NewUnreliableSimulatedNode(netsim)
  go netsim.Run()

  for i := 0; i < 10; i++ {
    packet := NewPacket(n1.Id(), n2.Id(), nil)
    n1.Send() <- packet

    recvPacket := (<-n2.Recv()).(Packet)
    assertEqual(t, packet.Src(), recvPacket.Src())
    assertEqual(t, packet.Dest(), recvPacket.Dest())
    assertEqual(t, packet.Payload(), recvPacket.Payload())
  }

  netsim.Stop()
}

func TestUnreliableNodesSameNumberOfSentAndReceivedPackets(t *testing.T) {
  netsim := testUnderlayNetsim(100)

  nbrPackets := 50
  nbrNodes   := 50

  nodes := []UnreliableNode{}
  for i := 0; i < nbrNodes; i++ {
    nodes = append(nodes, NewUnreliableSimulatedNode(netsim))
  }
  go netsim.Run()

  for i := 0; i < nbrNodes; i++ {
    n1 := nodes[i]
    n2 := nodes[nbrNodes - i - 1]
    for j := 0; j < nbrPackets; j++ {
      packet := NewPacket(n1.Id(), n2.Id(), nil)
      n1.Send() <- packet
    }
  }
  for i := 0; i < nbrNodes; i++ {
    node := nodes[i]
    for j := 0; j < nbrPackets; j++ {
      <-node.Recv()
    }
    select {
    case <-node.Recv():
      t.Fatalf("More packets than expected arrived.")
    default:
    }
  }
}

func TestUnreliableNodesJoinReturnDifferentID(t *testing.T) {
  netsim := testUnderlayNetsim(10)

  // Join can be called only on at least 2 nodes network
  n1 := NewUnreliableSimulatedNode(netsim)
  n2 := NewUnreliableSimulatedNode(netsim)
  for i := 0; i < 10; i++ {
    assertNotEqual(t, n1.Id(), n1.Join())
    assertNotEqual(t, n2.Id(), n2.Join())
  }
}

func TestGetBootstrapResturnsSameBootstrap(t *testing.T) {
  netsim := testUnderlayNetsim(10)

  assertEqual(t, GetBootstrap(netsim), GetBootstrap(netsim))
}
