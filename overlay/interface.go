package overlay

import (
	"github.com/danalex97/Speer/events"
	"github.com/danalex97/Speer/interfaces"
	"github.com/danalex97/Speer/underlay"
)

var activeSet = make(map[*underlay.NetworkSimulation]LatencyMap)

type UnreliableNode interface {
	events.Decorable

	interfaces.UnreliableNode
}

// An UnreliableSimulatedNode uses a LatencyConnector along with a Bootstrap to allow
// interaction with the simulator by providing an utility to the user.
type UnreliableSimulatedNode struct {
	events.Decorable

	simulation *underlay.NetworkSimulation
	bridge     LatencyConnector
	bootstrap  Bootstrap
	id         string
}


// The Bootstrap is associated directly with the simulation. All the nodes
// need to refer to the same bootstrap, so we use a global map to associate
// a NetworkSimulation with a Bootstrap.
func GetBootstrap(simulation *underlay.NetworkSimulation) Bootstrap {
	netMap := NewNetworkMap(simulation.Network())
	if mp, ok := activeSet[simulation]; ok {
		netMap = mp
	} else {
		activeSet[simulation] = netMap
	}

	return netMap
}

func NewUnreliableSimulatedNode(simulation *underlay.NetworkSimulation) UnreliableNode {
	node := new(UnreliableSimulatedNode)

	var netMap LatencyMap
	if mp, ok := activeSet[simulation]; ok {
		netMap = mp
	} else {
		netMap = NewNetworkMap(simulation.Network())
		activeSet[simulation] = netMap
	}

	node.id = netMap.NewId()
	node.bridge = NewUnderlayChan(node.id, simulation, netMap)
	node.bootstrap = netMap
	node.simulation = simulation

	// The actual decorable is at bridge level.
	// To allow direct interfacing, we create a tunnnel.
	node.Decorable = events.NewTunnel(node.bridge)

	return node
}

func (n *UnreliableSimulatedNode) Id() string {
	return n.id
}

func (n *UnreliableSimulatedNode) Send(msg interface{}) {
	// n.bridge.Send(msg)
}

func (n *UnreliableSimulatedNode) Recv() <-chan interface{} {
	// return n.bridge.Recv()
	return nil
}

func (n *UnreliableSimulatedNode) Join() string {
	return n.bootstrap.Join(n.id)
}

func (n *UnreliableSimulatedNode) Simulation() *underlay.NetworkSimulation {
	return n.simulation
}
