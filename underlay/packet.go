package underlay

type Packet interface {
	Src() Router
	Dest() Router
	Payload() interface{}
}

type packet struct {
	src     Router
	dest    Router
	payload interface{}
}

func NewPacket(src, dest Router, payload interface{}) Packet {
	pkt := new(packet)
	pkt.src = src
	pkt.dest = dest
	pkt.payload = payload
	return pkt
}

func (p *packet) Src() Router {
	return p.src
}

func (p *packet) Dest() Router {
	return p.dest
}

func (p *packet) Payload() interface{} {
	return p.payload
}
