package overlay

import (
  "github.com/danalex97/Speer/underlay"
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
  simulation *underlay.NetworkSimulation
  bootstrap   Bootstrap
}

func NewUnderlayChan(id string, simulation *underlay.NetworkSimulation, bootstrap Bootstrap) Chan {
  chn := new(UnderlayChan)

  chn.id = id
  chn.simulation = simulation
  chn.bootstrap  = bootstrap

  chn.send = make(chan interface{})
  chn.recv = make(chan interface{})

  go chn.establishListeners()
  go chn.establishPushers()

  return chn
}

func (u *UnderlayChan) establishListeners() {
  obs := NewEventObserver(u.bootstrap.Router(u.id))
  u.simulation.RegisterObserver(obs)

  for {
    event  := <- obs.EventChan()
    packet := event.Payload().(underlay.Packet)
    u.recv <- packet
  }
}

func (u *UnderlayChan) establishPushers() {
  for {
    // event := <- u.send
  }
}

func (u *UnderlayChan) Send() chan<- interface{} {
  return u.send
}

func (u *UnderlayChan) Recv() <-chan interface{} {
  return u.recv
}

func (u *UnderlayChan) UnderlayPacket(p *Packet) *underlay.Packet {
  return underlay.NewPacket(
    u.bootstrap.Router(p.Src()),
    u.bootstrap.Router(p.Dest()),
    p.Payload(),
  )
}

func (u *UnderlayChan) OverlayPacket(p *underlay.Packet) *Packet {
  return NewPacket(
    u.bootstrap.Id(p.Src()),
    u.bootstrap.Id(p.Dest()),
    p.Payload(),
  )
}
