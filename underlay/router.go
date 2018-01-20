package underlay

import (
  "error"
  . "github.com/danalex97/Speer/events"
)

type Router interface {
  Receiver
  Connect(Connection) error
}

type shortestPathRouter struct {
  table []Connection
}

type sprPacket struct {
  packet
  path []Connection
}

func NewShortestPathRouter() Router {
  router := shortestPathRouter{[]shortestPathRouter{}}
  return router
}

func (r *shortestPathRouter) Connect(conn Connection) {
  latency := conn.Latency()
  router  := conn.Router()

  spr, ok := router.(shortestPathRouter)
  if !ok {
    return error("Connection type not supported.")
  }

  *table = append(*table, conn)

  return nil
}

func (r *shortestPathRouter) Receive(event *Event) {
  switch payload := event.Payload().(type) {
  case packet:
    // first hop
    nextPayload, ok := dijkstra(r, payload.dest)
    if !ok {
      return nil
    }
    return buildNextEvent(nextPayload)

  case sprPacket:
    // next hops
    nextPayload := payload
    if len(nextPayload.path) == 0 {
      return nil
    }
    return buildNextEvent(nextPayload)
  }
  return nil
}

func buildNextEvent(nextPayload sprPacket) *Event {
  conn := nextPayload.path[0]
  nextPayload.path = nextPayload.path[1:]
  return NewEvent(
    event.Timestamp() + conn.Latency(),
    nextPayload,
    conn.Router()
  )
}

func dijkstra(r *shortestPathRouter) sprPacket, bool {
}
