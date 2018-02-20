package overlay

import (
  "testing"
  "math"
  "math/rand"
  "time"
  "github.com/danalex97/Speer/underlay"
)

func assertEqual(t *testing.T, a, b interface {}) {
  if a != b {
    t.Fatalf("%s == %s", a, b)
  }
}

func assertNotEqual(t *testing.T, a, b interface {}) {
  if a == b {
    t.Fatalf("%s == %s", a, b)
  }
}

func testNetworkMap(nodes int) OverlayMap {
  rand.Seed(time.Now().UTC().UnixNano())
  edges := int(math.Log2(float64(nodes)) * float64(nodes) / 2 + 1)
  network := underlay.NewRandomUniformNetwork(nodes, edges, 2, 10)
  return NewNetworkMap(network)
}

func TestNetworkMapNewIdReturnsDifferentRouter(t *testing.T) {
  netmap := testNetworkMap(10)
  ids := []string{}
  for i := 0; i < 10; i++ {
    ids = append(ids, netmap.NewId())
  }
  for i := 0; i < 10; i++ {
    for j := 0; j < i; j++ {
      assertNotEqual(t, ids[i], ids[j])
    }
  }
}

func TestNetworkMapJoinReturnsDifferentRouter(t *testing.T) {
  netmap := testNetworkMap(10)
  for i := 0; i < 10; i++ {
    netmap.NewId()
  }
  for i := 0; i < 10; i++ {
    j := netmap.Join(string(i))
    assertNotEqual(t, string(i), j)
  }
}

func TestNetworkMapCanAccessRouter(t *testing.T) {
  netmap := testNetworkMap(10)
  ids := []string{}
  for i := 0; i < 10; i++ {
    ids = append(ids, netmap.NewId())
  }

  for i := 0; i < 10; i++ {
    for j := 0; j < i; j++ {
      assertNotEqual(t, netmap.Router(ids[i]), netmap.Router(ids[j]))
    }
  }
}

func TestNetworkMapCanAccessId(t *testing.T) {
  netmap := testNetworkMap(10)
  ids := []string{}
  for i := 0; i < 10; i++ {
    ids = append(ids, netmap.NewId())
  }

  for i := 0; i < 10; i++ {
    assertEqual(t, netmap.Id(netmap.Router(ids[i])), ids[i])
  }
}
