package sdk

import (
	"github.com/danalex97/Speer/events"
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

	simulation *events.Simulation

	time      func() int
	bootstrap overlay.Bootstrap
	id        string
}

func NewSimulatedNode(
	controlTransport interfaces.ControlTransport,
	dataTransport interfaces.DataTransport,
	simulation *events.Simulation,
	bootstrap overlay.Bootstrap,
	id string,
	time func() int,
) NodeUtil {
	return &SimulatedNode{
		Connector: &Connector{
			ControlTransport: controlTransport,
			DataTransport:    dataTransport,
		},
		simulation: simulation,
		time:       time,
		bootstrap:  bootstrap,
		id:         id,
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

func (n *SimulatedNode) Routine(interval int, routine func()) interfaces.Routine {
	routineReceiver := events.NewRoutine(interval, routine)
	n.simulation.Push(events.NewEvent(
		n.simulation.Time(),
		nil,
		routineReceiver,
	))
	return routineReceiver
}

func (n *SimulatedNode) Callback(timeout int, routine func()) interfaces.Callback {
	callbackReceiver := events.NewCallback(routine)
	n.simulation.Push(events.NewEvent(
		n.simulation.Time()+timeout,
		nil,
		callbackReceiver,
	))
	return callbackReceiver
}
