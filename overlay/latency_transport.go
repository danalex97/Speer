package overlay

import (
	"github.com/danalex97/Speer/events"
	"github.com/danalex97/Speer/interfaces"
)

type ControlTransport interface {
	events.Decorable
	interfaces.ControlTransport
}

type SimulatedNetworkedTransport struct {
	events.Decorable

	// overlay node id
	id string

	// the latency connector is used to send and receive data
	connector LatencyConnector
}

func NewSimulatedNetworkedTransport(
	id string,
	connector LatencyConnector,
) ControlTransport {
	transport := &SimulatedNetworkedTransport{
		Decorable : events.NewTunnel(connector),

		id : id,
		connector : connector,
	}

	transport.SetProxy(stripPayload)

	return transport
}

func (t *SimulatedNetworkedTransport) ControlSend(
	dst string,
	msg interface{},
) {
	t.connector.Send(NewPacket(t.id, dst, msg))
}

func stripPayload(m interface {}) interface {} {
  return m.(Packet).Payload()
}

func (t *SimulatedNetworkedTransport) ControlRecv() <-chan interface{} {
  return t.connector.Recv()
}

// this is kept only for backwards compatibility, we want to remove it
func (t *SimulatedNetworkedTransport) ControlPing(id string) bool {
	return true
}
