package events

type Proxy func(interface {}) interface{}

func Identity(m interface {}) interface{} {
  return m
}

type Decorable interface {
  GetProxy() Proxy
  SetProxy(Proxy)
}

type Decorator struct {
  Proxy
}

func NewDecorator() Decorable {
  return &Decorator{Identity}
}

func (d *Decorator) GetProxy() Proxy {
  return d.Proxy
}

func (d *Decorator) SetProxy(p Proxy) {
  d.Proxy = p
}
