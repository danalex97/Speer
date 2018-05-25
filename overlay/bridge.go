package overlay

import (
  . "github.com/danalex97/Speer/events"
  "github.com/danalex97/Speer/underlay"
)

type Bridge interface {
  Decorable

  Send(interface{})
  Recv() <-chan interface{}
}

type UnderlayChan struct {
  *Decorator

  id string

  simulation *underlay.NetworkSimulation
  netMap     OverlayMap

  observer   DecorableObserver
}

func NewUnderlayChan(
    id         string,
    simulation *underlay.NetworkSimulation,
    netMap     OverlayMap) Bridge {

  u := new(UnderlayChan)

  u.id         = id
  u.simulation = simulation
  u.netMap     = netMap

  // Allow decoration at bigger levels.
  u.Decorator = NewDecorator()

  // Establish listener
  u.observer = NewEventObserver(u.netMap.Router(u.id))
  u.observer.SetProxy(u.ReceiveEvent)
  u.simulation.RegisterObserver(u.observer)

  return u
}

// Notify observer directly by creating an event and delivering it to the
// observer directly.
func (u *UnderlayChan) notifyPacket(packet underlay.Packet) {
  // We need to run this in a separate routine since enqueing can be blocking,
  // resulting in a problem when sending a packet to self.
  go u.observer.EnqueEvent(NewEvent(0, packet, packet.Dest()))
}

func (u *UnderlayChan) ReceiveEvent(m interface {}) interface{} {
  event  := (m).(*Event)
  packet := event.Payload().(underlay.Packet)
  overPacket := u.OverlayPacket(packet)

  if packet.Src() == nil {
    return nil
  }

  // We need to look only at our own packets.
  if overPacket.Dest() != u.id {
    return nil
  }
  // fmt.Printf("Packet delivered: {%s, %s}\n", overPacket.Src(), overPacket.Dest())

  return u.Proxy(overPacket)
}

func (u *UnderlayChan) Send(msg interface {}) {
  overPacket := msg.(Packet)
  packet     := u.UnderlayPacket(overPacket)

  if u.id == overPacket.Dest() {
    // Packet sent to self.
    u.notifyPacket(packet)
    return
  }

  u.simulation.SendPacket(packet)
}

func (u *UnderlayChan) Recv() <-chan interface{} {
  return u.observer.Recv()
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
