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
}
