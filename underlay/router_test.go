package underlay

import (
  . "github.com/danalex97/Speer/events"
  "testing"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

func TestShortestPathRouterRingTopology(t *testing.T) {
  routers := make([]Router, 0)
  n := 50
  for i := 0; i < n; i++ {
    routers = append(routers, NewShortestPathRouter())
  }
  for i := 0; i < n; i++ {
    routers[i].Connect(NewStaticConnection(1, routers[(i + 1) % n]))
    routers[i].Connect(NewStaticConnection(1, routers[(i - 1 + n) % n]))
  }

  for i := 0; i < n; i++ {
    pkt := NewPacket(routers[i], routers[(i + 5) % n], nil)
    pkt2 := NewPacket(routers[i], routers[(i - 5 + n) % n], nil)

    e := NewEvent(0, pkt, routers[i])
    e2 := NewEvent(0, pkt2, routers[i])

    assertEqual(t, routers[i].Receive(e).Receiver(), routers[(i + 1) % n])
    assertEqual(t, routers[i].Receive(e2).Receiver(), routers[(i - 1 + n) % n])
  }
}
