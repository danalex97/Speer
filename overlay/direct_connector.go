package overlay

import (
	"github.com/danalex97/Speer/interfaces"
)

const controlMessageCapacity int = 1000000

type DirectConnector interface {
	interfaces.ControlTransport

	Chan() chan interface{}
}

// An DirectChan implements a DirectConnector by using channels directly.
type DirectChan struct {
	id string
	recv chan interface {}

	networkMap DirectMap
}

func NewDirectChan(
	id string,
	networkMap DirectMap,
) DirectConnector {
	return &DirectChan{
		id: id,
		recv: make(chan interface {}, controlMessageCapacity),

		networkMap: networkMap,
	}
}

func (d *DirectChan) ControlSend(dst string, msg interface{}) {
	directChan := d.networkMap.Chan(dst)
	if directChan != nil {
		directChan.Chan() <- msg
	}
}

func (d *DirectChan) ControlRecv() <-chan interface{} {
	return d.recv
}

func (d *DirectChan) ControlPing(id string) bool {
	return d.networkMap.Chan(id) != nil
}

func (d *DirectChan) Chan() chan interface{} {
	return d.recv
}
