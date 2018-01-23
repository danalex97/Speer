package overlay

// should work with channs (?)

type Chan interface {
  Send() chan<- interface{}
  Recv() <-chan interface{}
}

type UnderlayChan struct {
  id string
  send chan<- interface{}
  recv <-chan interface{}
  simulation *NetworkSimulation
  bootstrap   Bootstrap
}

func NewUnderlayChan(id string, simulation *NetworkSimulation, bootstrap Bootstrap) Chan {
  chn := new(UnderlayChan)

  chn.send = make(chan<- interface{})
  chn.recv = make(<-chan interface{})

  chn.id = id
  chn.simulation = simulation
  chn.bootstrap  = bootstrap

  go chn.establishListeners()

  return chn
}

func (u *Underlay) establishListeners() {
  // need to use observer pattern over the simulation 
}

func (u *Underlay) Send() chan<- interface{} {
  return u.send
}

func (u *Underlay) Recv() <-chan interface{} {
  return u.recv
}
