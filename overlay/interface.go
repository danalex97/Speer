package overlay

import (
  . "github.com/danalex97/Speer/underlay"
  . "github.com/danalex97/Speer/events"
)

type Chan interface {
  Send() chan<- interface{}
  Recv() <-chan interface{}
}

type UnderlayChan struct {
  id string
  send chan interface{}
  recv chan interface{}
  simulation *NetworkSimulation
  bootstrap   Bootstrap
}

func NewUnderlayChan(id string, simulation *NetworkSimulation, bootstrap Bootstrap) Chan {
  chn := new(UnderlayChan)

  chn.id = id
  chn.simulation = simulation
  chn.bootstrap  = bootstrap

  chn.send = make(chan interface{})
  chn.recv = make(chan interface{})

  go chn.establishListeners()

  return chn
}

func (u *UnderlayChan) establishListeners() {
  obs := NewEventObserver(u.bootstrap.Router(u.id))
  u.simulation.RegisterObserver(obs)

  for {
    event := <- obs.EventChan()
    u.recv <- event.Payload()
  }
}

func (u *UnderlayChan) Send() chan<- interface{} {
  return u.send
}

func (u *UnderlayChan) Recv() <-chan interface{} {
  return u.recv
}
