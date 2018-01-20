package underlay

type Connection interface {
  Latency() int
  Router() Router
}

type staticConnection struct {
  latency int
  router  Router
}

func NewStaticConnection(latency int, router Router) Connection {
  conn := new(staticConnection)
  conn.latency = latency
  conn.router  = router
  return conn
}

func (conn staticConnection) Latency() {
  return conn.latency
}

func (conn staticConnection) Router() {
  return conn.router
}
