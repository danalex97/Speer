package underlay

import (
  . "github.com/danalex97/Speer/events"
)

type Router interface {
  Receiver
  Connect(Connection) error
  Connections() []Connection
}

type sprPacket struct {
  packet
  path []Connection
}

type shortestPathRouter struct {
  table []Connection
}

func NewShortestPathRouter() Router {
  router := new(shortestPathRouter)
  router.table = []Connection{}
  return router
}

func (r *shortestPathRouter) Connect(conn Connection) error {
  r.table = append(r.table, conn)

  return nil
}

func (r *shortestPathRouter) Connections() []Connection {
  return r.table;
}

func (r *shortestPathRouter) Receive(event *Event) *Event {
  switch payload := event.Payload().(type) {
  case *packet:
    // first or last hop
    nextPayload, ok := bellman(payload, r, payload.Dest())
    if !ok {
      return nil
    }
    return buildNextEvent(event, nextPayload)

  case *sprPacket:
    // next hops
    nextPayload := payload
    if len(nextPayload.path) == 0 {
      // the packet arrived at destination so don't to anything
      // the packet is used by the observers
      return nil
    }
    return buildNextEvent(event, nextPayload)
  }
  return nil
}

func buildNextEvent(event *Event, nextPayload *sprPacket) *Event {
  conn := nextPayload.path[0]
  nextPayload.path = nextPayload.path[1:]

  return NewEvent(
    event.Timestamp() + conn.Latency(),
    nextPayload,
    conn.Router(),
  )
}

// not very good, it does not clear doubled values as well
// used only as a policy for small tests
func bellman(packet Packet, src *shortestPathRouter, dest Router) (*sprPacket, bool) {
  eq := NewLazyEventQueue()
  eq.Push(NewEvent(0, 0, nil))

  conns := []Connection{NewStaticConnection(0, src)}
  last  := []int{-1}
  ctr   := 0

  pkt := new(sprPacket)
  pkt.src  = packet.Src()
  pkt.dest = packet.Dest()
  pkt.payload  = packet.Payload()

  for {
    curr := eq.Pop()
    if curr == nil {
      return pkt, false
    }

    cost := curr.Timestamp()
    idx  := curr.Payload().(int)
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

    for i := range(router.table) {
      conn := router.table[i]

      ctr += 1
      conns = append(conns, conn)
      last  = append(last, idx)

      eq.Push(NewEvent(cost + conn.Latency(), ctr, nil))
    }
  }
}
