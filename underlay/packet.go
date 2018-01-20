package underlay

type Packet struct {
  src  Router
  dest Router
}

func NewPacket(src, dest Router) *Packet {
  pkt := new(Packet)
  pkt.src = src
  pkt.dest = dest
  return pkt
}

func (p *Packet) Src() Router {
  return p.src
}

func (p *Packet) Dest() Router {
  return p.dest
}
