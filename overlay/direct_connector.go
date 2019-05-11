package overlay

import (
	"github.com/danalex97/Speer/interfaces"
)

const controlMessageCapacity int = 1000000

type DirectConnector interface {
	interfaces.ControlTransport
}

// An DirectChan implements a DirectConnector by using channels directly.
type DirectChan struct {
	id string
	recv chan interface {}

	networkMap LatencyMap
}

func NewDirectChan(
	id string,
	networkMap LatencyMap,
) DirectConnector {
	return &DirectChan{
		id: id,
		recv: make(chan interface {}, controlMessageCapacity),

		networkMap: networkMap,
	}
}

func (d *DirectChan) ControlSend(dst string, msg interface{}) {
}

func (d *DirectChan) ControlRecv() <-chan interface{} {
	return d.recv
}

func (u *DirectChan) ControlPing(id string) bool {
	return true
}
