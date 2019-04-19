package underlay

import (
	"fmt"
	. "github.com/danalex97/Speer/events"
)

const RouterCacheSize int = 50

// A Router is a Receiver which schedules new events at other routers according
// to a routing algorithm. We use the Bellman-Ford alogrithm to determine the
// shortest paths.
type Router interface {
	Receiver

	Connect(Connection) error
	Connections() []Connection

	Domain() string
}

type sprPacket struct {
	packet
	path []Connection
}

type shortestPathRouter struct {
	table  []Connection
	cache  Cache
	domain string
}

func NewShortestPathRouter(domain string) Router {
	router := new(shortestPathRouter)
	router.table = []Connection{}
	router.cache, _ = NewLRUCache(RouterCacheSize)
	router.domain = domain
	return router
}

func (r *shortestPathRouter) Domain() string {
	return r.domain
}

func (r *shortestPathRouter) Connect(conn Connection) error {
	r.table = append(r.table, conn)

	return nil
}

func (r *shortestPathRouter) Connections() []Connection {
	return r.table
}

func (r *shortestPathRouter) logPath(payload *sprPacket) {
	fmt.Printf("%p - Road payload: %v [", r, payload.packet)
	for i := 0; i < len(payload.path); i++ {
		fmt.Printf(" %v", payload.path[i])
	}
	fmt.Printf(" ]\n")
}

func (r *shortestPathRouter) Receive(event *Event) *Event {
	switch payload := event.Payload().(type) {
	case *packet:
		// check cache
		if el, ok := r.cache.Get(payload.dest); ok {
			conn := el.(Connection)
			return buildNextCachedEvent(event, conn, payload)
		}
		// fmt.Printf("%p - Original payload: %v\n", r, payload)

		// first or last hop
		nextPayload, ok := bellman(payload, r, payload.Dest())
		//r.logPath(nextPayload)

		if !ok || len(nextPayload.path) == 0 {
			return nil
		}
		return r.buildNextEvent(event, nextPayload)

	case *sprPacket:
		// next hops
		nextPayload := payload
		//r.logPath(payload)

		if len(nextPayload.path) == 0 {
			// the packet arrived at destination so don't to anything
			// the packet is used by the observers
			return nil
		}
		return r.buildNextEvent(event, nextPayload)
	}
	return nil
}

func buildNextCachedEvent(event *Event, conn Connection, payload interface{}) *Event {
	return NewEvent(
		event.Timestamp()+conn.Latency(),
		payload,
		conn.Router(),
	)
}

func (r *shortestPathRouter) buildNextEvent(event *Event, nextPayload *sprPacket) *Event {
	// next hop
	conn := nextPayload.path[0]
	nextPayload.path = nextPayload.path[1:]

	// cache the path
	r.cache.Put(nextPayload.dest, conn)

	// return the corresponding event
	return NewEvent(
		event.Timestamp()+conn.Latency(),
		nextPayload,
		conn.Router(),
	)
}

// We replace the Dijkstra’s algorithm with a Bellman-Ford variation, adding
// the path that the packet will follow as its payload. This results in a memory
// overhead of only O(E) with time overhead of O(V*E), with a O(E) complexity
// in the average case.
func bellman(packet Packet, src *shortestPathRouter, dest Router) (*sprPacket, bool) {
	eq := NewLazyEventQueue()
	eq.Push(NewEvent(0, 0, nil))

	conns := []Connection{NewStaticConnection(0, src)}
	last := []int{-1}
	ctr := 0
	best := make(map[Router]int)

	pkt := new(sprPacket)
	pkt.src = packet.Src()
	pkt.dest = packet.Dest()
	pkt.payload = packet.Payload()

	best[pkt.src] = 0
	for {
		curr := eq.Pop()
		if curr == nil {
			return pkt, false
		}

		cost := curr.Timestamp()
		idx := curr.Payload().(int)
		router := conns[idx].Router().(*shortestPathRouter)

		if router == dest {
			// build the packet
			idxs := []int{}
			for i := idx; i > 0; i = last[i] {
				idxs = append(idxs, i)
			}
			idxs = append(idxs, 0)

			pkt.path = []Connection{}
			for i := len(idxs) - 2; i >= 0; i -= 1 {
				pkt.path = append(pkt.path, conns[idxs[i]])
			}
			return pkt, true
		}

		for i := range router.table {
			conn := router.table[i]
			r := conn.Router()

			if cst, ok := best[r]; !ok || cst > cost+conn.Latency() {
				best[r] = cost + conn.Latency()
			}

			if cost+conn.Latency() > best[r] {
				continue
			}

			ctr += 1
			conns = append(conns, conn)
			last = append(last, idx)

			eq.Push(NewEvent(cost+conn.Latency(), ctr, nil))
		}
	}
}
