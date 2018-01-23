package underlay

type Packet struct {
  src     Router
  dest    Router
  payload interface{}
}

func NewPacket(src, dest Router, payload interface{}) *Packet {
  pkt := new(Packet)
  pkt.src = src
  pkt.dest = dest
  pkt.payload = payload
  return pkt
}

func (p *Packet) Src() Router {
  return p.src
}

func (p *Packet) Dest() Router {
  return p.dest
}

func (p *Packet) Payload() interface{} {
  return p.payload
}
