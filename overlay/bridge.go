package overlay

import (
  . "github.com/danalex97/Speer/events"
  "github.com/danalex97/Speer/underlay"
  // "runtime"
  "fmt"
)

type Bridge interface {
  Send(interface{})
  Recv() <-chan interface{}
}

type UnderlayChan struct {
  id string

  send chan interface{}
  recv chan interface{}

  simulation *underlay.NetworkSimulation
  netMap     OverlayMap
}

const sendSize int = 50
const recvSize int = 10000

func NewUnderlayChan(
    id         string,
    simulation *underlay.NetworkSimulation,
    netMap     OverlayMap) Bridge {

  chn := new(UnderlayChan)

  chn.id         = id
  chn.simulation = simulation
  chn.netMap     = netMap

  chn.send = make(chan interface{}, sendSize)
  chn.recv = make(chan interface{}, recvSize)

  go chn.establishListeners()

  return chn
}

func (u *UnderlayChan) notifyRecvPkt(overPacket Packet) {
  select {
  case u.recv <- overPacket:
  default:
    // Packet dropped when receiver queue is full
    fmt.Println("Receiver queue full, packet dropped!")
  }
}

func (u *UnderlayChan) establishListeners() {
  obs := NewEventObserver(u.netMap.Router(u.id))
  u.simulation.RegisterObserver(obs)

  for {
    event := (<-obs.Recv()).(*Event)
    packet := event.Payload().(underlay.Packet)
    overPacket := u.OverlayPacket(packet)

    if packet.Src() == nil {
      continue
    }
    if overPacket.Src() == u.id {
      continue
    }

    // We need to look only at our own packets.
    if overPacket.Dest() != u.id {
      continue
    }
    // fmt.Printf("Packet delivered: {%s, %s}\n", overPacket.Src(), overPacket.Dest())

    u.notifyRecvPkt(overPacket)
  }
}

func (u *UnderlayChan) Send(msg interface {}) {
  overPacket := msg.(Packet)
  if u.id == overPacket.Dest() {
    // Packet sent to self.
    u.notifyRecvPkt(overPacket)
    return
  }

  packet := u.UnderlayPacket(overPacket)
  u.simulation.SendPacket(packet)
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
