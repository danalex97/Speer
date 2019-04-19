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

func TestSimpleTopology(t *testing.T) {
	r1 := NewShortestPathRouter("")
	r2 := NewShortestPathRouter("")
	r3 := NewShortestPathRouter("")

	r1.Connect(NewStaticConnection(1, r2))
	r2.Connect(NewStaticConnection(2, r3))
	r1.Connect(NewStaticConnection(10, r3))

	pkt := NewPacket(r1, r3, nil)
	e := NewEvent(0, pkt, r1)
	e2 := r1.Receive(e)
	assertEqual(t, e2.Receiver(), r2)
	assertEqual(t, e2.Timestamp(), 1)
	e3 := r2.Receive(e2)
	assertEqual(t, e3.Receiver(), r3)
	assertEqual(t, e3.Timestamp(), 3)

	// The last step returns null and the packet does not need to get stripped
	var nilE *Event
	assertEqual(t, r3.Receive(e3), nilE)
}

func TestShortestPathRouterRingTopology(t *testing.T) {
	routers := make([]Router, 0)
	n := 50
	for i := 0; i < n; i++ {
		routers = append(routers, NewShortestPathRouter(""))
	}
	for i := 0; i < n; i++ {
		routers[i].Connect(NewStaticConnection(1, routers[(i+1)%n]))
		routers[i].Connect(NewStaticConnection(1, routers[(i-1+n)%n]))
	}

	for i := 0; i < n; i++ {
		pkt := NewPacket(routers[i], routers[(i+5)%n], nil)
		pkt2 := NewPacket(routers[i], routers[(i-5+n)%n], nil)

		e := NewEvent(0, pkt, routers[i])
		e2 := NewEvent(0, pkt2, routers[i])

		assertEqual(t, routers[i].Receive(e).Receiver(), routers[(i+1)%n])
		assertEqual(t, routers[i].Receive(e2).Receiver(), routers[(i-1+n)%n])
	}
}

func TestShortestPathRouterInternetwork(t *testing.T) {
	network := NewInternetwork(5, 10, 0, 0)
	routers := network.Routers
	n := len(routers)

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i != j {
				pkt := NewPacket(routers[i], routers[j], nil)
				e := NewEvent(0, pkt, routers[j])

				routers[i].Receive(e)
			}
		}
	}
}
