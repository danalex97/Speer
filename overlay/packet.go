package overlay

type Packet struct {
  src     string
  dest    string
  payload interface{}
}

func NewPacket(src, dest string, payload interface{}) *Packet {
  pkt := new(Packet)
  pkt.src = src
  pkt.dest = dest
  pkt.payload = payload
  return pkt
}

func (p *Packet) Src() string {
  return p.src
}

func (p *Packet) Dest() string {
  return p.dest
}

func (p *Packet) Payload() interface{} {
  return p.payload
}
