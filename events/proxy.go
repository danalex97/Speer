package events

type Proxy func(interface {}) interface{}

func Identity(m interface {}) interface{} {
  return m
}

type Decorable interface {
  SetProxy(Proxy)
}

type Decorator struct {
  Proxy
}

func NewDecorator() *Decorator {
  return &Decorator{Identity}
}

func (d *Decorator) SetProxy(p Proxy) {
  d.Proxy = p
}

type Tunnel struct {
  Decorable
}

func NewTunnel(d Decorable) *Tunnel {
  return &Tunnel{d}
}

func (t *Tunnel) SetProxy(p Proxy) {
  t.Decorable.SetProxy(p)
}
