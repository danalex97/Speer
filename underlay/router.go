package underlay

import (
  // "errors"
  // "fmt"
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
  // spr, ok := router.(shortestPathRouter)
  // if !ok {
  //   return errors.New("Connection type not supported.")
  // }

  r.table = append(r.table, conn)

  return nil
}

func (r *shortestPathRouter) Connections() []Connection {
  return r.table;
}

func (r *shortestPathRouter) Receive(event *Event) *Event {
  switch payload := event.Payload().(type) {
  case *packet:
    // first hop
    nextPayload, ok := bellman(r, payload.Dest())
    if !ok {
      return nil
    }
    return buildNextEvent(event, nextPayload)

  case *sprPacket:
    // next hops
    nextPayload := payload
    if len(nextPayload.path) == 0 {
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
func bellman(src *shortestPathRouter, dest Router) (*sprPacket, bool) {
  eq := NewLazyEventQueue()
  eq.Push(NewEvent(0, 0, nil))

  conns := []Connection{NewStaticConnection(0, src)}
  last  := []int{-1}
  ctr   := 0

  for {
    curr := eq.Pop()
    if curr == nil {
      return new(sprPacket), false
    }

    cost := curr.Timestamp()
    idx  := curr.Payload().(int)
    router := conns[idx].Router().(*shortestPathRouter)

    for i := range(router.table) {
      conn := router.table[i]

      ctr += 1
      conns = append(conns, conn)
      last  = append(last, idx)

      eq.Push(NewEvent(cost + conn.Latency(), ctr, nil))
      if conn.Router() == dest {
        // build the packet
        idxs := []int{}
        for i := ctr; i > 0; i = last[i] {
          idxs = append(idxs, ctr)
        }
        idxs = append(idxs, 0)

        pkt := new(sprPacket)
        pkt.path = []Connection{}
        for i := len(idxs) - 1; i >= 0; i -= 1 {
          pkt.path = append(pkt.path, conns[i])
        }
        return pkt, true
      }
    }
  }
}
