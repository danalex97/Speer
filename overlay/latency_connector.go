package overlay

import (
	. "github.com/danalex97/Speer/events"
	"github.com/danalex97/Speer/interfaces"
	"github.com/danalex97/Speer/underlay"
)

// A LatencyConnector is a decorable interface which allows sending and
// receiving packets.
type LatencyConnector interface {
	Decorable

	interfaces.ControlTransport
}

// The UnderlayChan implements a LatencyConnector by using a proxy to strip
// the payload of the underlay packets at receiving a packet and sending a
// packet is done by using the network map to decorate the overlay packet
// inside an underlay packet.
//
// The mechanism used for delivering packets is an observer attached to the
// router corresponing to the overlay id.
type UnderlayChan struct {
	*Decorator

	id string

	simulation *underlay.NetworkSimulation
	networkMap LatencyMap

	observer DecorableObserver
}

func NewUnderlayChan(
	id string,
	simulation *underlay.NetworkSimulation,
	networkMap LatencyMap,
) LatencyConnector {
	u := new(UnderlayChan)

	u.id = id
	u.simulation = simulation
	u.networkMap = networkMap

	// Allow decoration at bigger levels.
	u.Decorator = NewDecorator()

	// Establish listener
	u.observer = NewEventObserver(u.networkMap.Router(u.id))
	u.observer.SetProxy(u.ReceiveEvent)
	u.simulation.RegisterObserver(u.observer)

	u.SetProxy(stripPayload)

	return u
}

// Notify observer directly by creating an event and delivering it to the
// observer directly.
func (u *UnderlayChan) notifyPacket(packet underlay.Packet) {
	// We need to run this in a separate routine since enqueing can be blocking,
	// resulting in a problem when sending a packet to self.
	go u.observer.EnqueEvent(NewEvent(0, packet, packet.Dest()))
}

// Proxy function used to strip the contents of an underlay packet. The
// UnderlayChan chan is a Decorator, so we call the Proxy function before
// delivering the packet.
func (u *UnderlayChan) ReceiveEvent(m interface{}) interface{} {
	event := (m).(*Event)
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

func (u *UnderlayChan) OverlayPacket(p underlay.Packet) Packet {
	return NewPacket(
		u.networkMap.Id(p.Src()),
		u.networkMap.Id(p.Dest()),
		p.Payload(),
	)
}

func (u *UnderlayChan) ControlSend(dst string, msg interface{}) {
	packet := underlay.NewPacket(
		u.networkMap.Router(u.id),
		u.networkMap.Router(dst),
		msg,
	)

	if u.id == dst {
		// Packet sent to self.
		u.notifyPacket(packet)
		return
	}

	u.simulation.SendPacket(packet)
}

func stripPayload(m interface{}) interface{} {
	return m.(Packet).Payload()
}

func (u *UnderlayChan) ControlRecv() <-chan interface{} {
	return u.observer.Recv()
}

func (u *UnderlayChan) ControlPing(id string) bool {
	return true
}
