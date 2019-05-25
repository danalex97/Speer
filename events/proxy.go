package events

// A Proxy is a function to be executed before pushing a message
// into a channel.
type Proxy func(interface{}) interface{}

func Identity(m interface{}) interface{} {
	return m
}

// An interface is a Decorable if it provides a proxy to be set.
// Before a message is pushed into a channel the proxy should
// process the message.
type Decorable interface {
	SetProxy(Proxy)
}

// A Decorator is the implementation of a Decorable.
type Decorator struct {
	Proxy
}

func NewDecorator() *Decorator {
	return &Decorator{Identity}
}

func (d *Decorator) SetProxy(p Proxy) {
	d.Proxy = p
}

// Creates a Proxy from a function that doesn't return an interface.
func NewProxy(f func()) Proxy {
	ret := func(_ interface{}) interface{} {
		f()
		return nil
	}
	return ret
}
