package sdk

import (
	"github.com/danalex97/Speer/interfaces"
	"github.com/danalex97/Speer/overlay"
)

type NodeUtil interface {
	interfaces.NodeUtil
}

type Connector struct {
	interfaces.ControlTransport
	interfaces.DataTransport
}

type SimulatedNode struct {
	*Connector

	time      func() int
	bootstrap overlay.Bootstrap
	id        string
}

func NewSimulatedNode(
	controlTransport interfaces.ControlTransport,
	dataTransport interfaces.DataTransport,
	bootstrap overlay.Bootstrap,
	id string,
	time func() int,
) NodeUtil {
	return &SimulatedNode{
		Connector: &Connector{
			ControlTransport: controlTransport,
			DataTransport:    dataTransport,
		},
		time:      time,
		bootstrap: bootstrap,
		id:        id,
	}
}

func (n *SimulatedNode) Id() string {
	return n.id
}

func (n *SimulatedNode) Join() string {
	return n.bootstrap.Join(n.id)
}

func (n *SimulatedNode) Time() func() int {
	return n.time
}

func (n *SimulatedNode) Transport() interfaces.Transport {
	return n.Connector
}
