package overlay

import (
	. "github.com/danalex97/Speer/events"
	"github.com/danalex97/Speer/interfaces"
	"github.com/danalex97/Speer/underlay"
)

// A LatencyConnector is an interface which allows sending and
// receiving packets. Moreover, the interface allows access to an
// ActiveObserver which can be used to react to packet receival.
type LatencyConnector interface {
	Observer() ActiveObserver

	interfaces.ControlTransport
}

// The UnderlayChan implements a LatencyConnector by using a proxy to strip
// the payload of the underlay packets at receiving a packet and sending a
// packet is done by using the network map to decorate the overlay packet
// inside an underlay packet.
//
// The mechanism used for delivering packets is a PassiveObserver attached to
// the router corresponing to the overlay id.
type UnderlayChan struct {
	id            string
	eventReceiver Receiver

	simulation *underlay.NetworkSimulation
	networkMap LatencyMap

	passiveObserver PassiveObserver
	activeObserver  ActiveObserver
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

	// Establish listener
	u.passiveObserver = NewPassiveEventObserver(u.networkMap.Router(u.id))
	u.passiveObserver.SetProxy(u.ReceiveEvent)

	// Create active observer
	u.activeObserver = NewActiveEventObserver(u.networkMap.Router(u.id))

	// Register observers; Note that the active observer should be attached
	// AFTER the passive observer so that the event delivered by the passive
	// observer in enqueued when the active observer while allow its Proxy to
	// be called. !Note this relies on details of the events package.!
	u.simulation.RegisterObserver(u.passiveObserver)
	u.simulation.RegisterObserver(u.activeObserver)

	return u
}

func (u *UnderlayChan) Observer() ActiveObserver {
	return u.activeObserver
}

// Proxy function used to strip the contents of an underlay packet. The
// UnderlayChan chan is a Decorator, so we call the Proxy function before
// delivering the packet.
func (u *UnderlayChan) ReceiveEvent(m interface{}) interface{} {
	event := (m).(*Event)
	packet := event.Payload().(underlay.Packet)
	overPacket := u.overlayPacket(packet)

	if packet.Src() == nil {
		return nil
	}

	// We need to look only at our own packets.
	if overPacket.Dest() != u.id {
		return nil
	}
	return stripPayload(overPacket)
}

func (u *UnderlayChan) ControlSend(dst string, msg interface{}) {
	packet := underlay.NewPacket(
		u.networkMap.Router(u.id),
		u.networkMap.Router(dst),
		msg,
	)

	if u.id == dst {
		// Packet sent to self.
		u.passiveObserver.Receive(NewEvent(0, packet, packet.Dest()))
		return
	}

	u.simulation.SendPacket(packet)
}

func (u *UnderlayChan) ControlRecv() <-chan interface{} {
	return u.passiveObserver.Recv()
}

func (u *UnderlayChan) ControlPing(id string) bool {
	return true
}

func stripPayload(m interface{}) interface{} {
	return m.(Packet).Payload()
}

func (u *UnderlayChan) overlayPacket(p underlay.Packet) Packet {
	return NewPacket(
		u.networkMap.Id(p.Src()),
		u.networkMap.Id(p.Dest()),
		p.Payload(),
	)
}
