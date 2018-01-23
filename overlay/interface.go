package overlay

// should work with channs (?)

type Chan interface {
  Send() chan<- interface{}
  Recv() <-chan interface{}
}

type UnderlayChan struct {
  send chan<- interface{}
  recv <-chan interface{}
  underlay *Network
}

func NewUnderlayChan(underlay *Network) Chan {
  chn := new(UnderlayChan)
  chn.send = make(chan<- interface{})
  chn.recv = make(<-chan interface{})
  chn.underlay = underlay

  chn.establishListeners()

  return chn
}

func (u *Underlay) establishListeners() {
}

func (u *Underlay) Send() chan<- interface{} {
  return u.send
}

func (u *Underlay) Recv() <-chan interface{} {
  return u.recv
}
