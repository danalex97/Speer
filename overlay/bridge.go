package overlay

import (
  "fmt"
  "github.com/danalex97/Speer/underlay"
  . "github.com/danalex97/Speer/events"
)

type Bridge interface {
  Send() chan<- interface{}
  Recv() <-chan interface{}
}

type UnderlayChan struct {
  id string
  send chan interface{}
  recv chan interface{}
  simulation *underlay.NetworkSimulation
  netMap     OverlayMap
}

func NewUnderlayChan(id string, simulation *underlay.NetworkSimulation, netMap OverlayMap) Bridge {
  chn := new(UnderlayChan)

  chn.id = id
  chn.simulation = simulation
  chn.netMap  = netMap

  chn.send = make(chan interface{})
  chn.recv = make(chan interface{})

  go chn.establishListeners()
  go chn.establishPushers()

  return chn
}

func (u *UnderlayChan) establishListeners() {
  obs := NewEventObserver(u.netMap.Router(u.id))
  u.simulation.RegisterObserver(obs)

  for {
    event  := <- obs.EventChan()
    packet := event.Payload().(underlay.Packet)
    overPacket := u.OverlayPacket(packet)

    if packet.Src() == nil {
      continue
    }
    if overPacket.Src() == u.id {
      continue
    }
    fmt.Printf("Packet delivered: {%s, %s}\n", overPacket.Src(), overPacket.Dest())

    u.recv <- overPacket
  }
}

func (u *UnderlayChan) establishPushers() {
  for {
    packet := u.UnderlayPacket((<- u.send).(Packet))
    u.simulation.SendPacket(packet)
  }
}

func (u *UnderlayChan) Send() chan<- interface{} {
  return u.send
}

func (u *UnderlayChan) Recv() <-chan interface{} {
  return u.recv
}

func (u *UnderlayChan) UnderlayPacket(p Packet) underlay.Packet {
  return underlay.NewPacket(
    u.netMap.Router(p.Src()),
    u.netMap.Router(p.Dest()),
    p.Payload(),
  )
}

func (u *UnderlayChan) OverlayPacket(p underlay.Packet) Packet {
  return NewPacket(
    u.netMap.Id(p.Src()),
    u.netMap.Id(p.Dest()),
    p.Payload(),
  )
}
