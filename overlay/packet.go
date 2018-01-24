package overlay

type Packet interface {
  Src()     string
  Dest()    string
  Payload() interface{}
}

type packet struct {
  src     string
  dest    string
  payload interface{}
}

func NewPacket(src, dest string, payload interface{}) Packet {
  pkt := new(packet)
  pkt.src = src
  pkt.dest = dest
  pkt.payload = payload
  return pkt
}

func (p *packet) Src() string {
  return p.src
}

func (p *packet) Dest() string {
  return p.dest
}

func (p *packet) Payload() interface{} {
  return p.payload
}
