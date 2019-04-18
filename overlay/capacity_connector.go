package overlay

import (
	. "github.com/danalex97/Speer/capacity"
	"github.com/danalex97/Speer/interfaces"
)

// Allow connecting between two nodes via scheduler
type CapacityConnector interface {
	interfaces.Node

	Connect(id string) interfaces.Link
}

type capacityConnector struct {
	capacityMap CapacityMap

	down int
	up   int
}

func NewCapacityConnector(
	up, down int,
	capacityMap CapacityMap,
) CapacityConnector {
	return &capacityConnector{
		capacityMap: capacityMap,

		down: down,
		up:   up,
	}
}

func (c *capacityConnector) Up() int {
	return c.up
}

func (c *capacityConnector) Down() int {
	return c.down
}

func (c *capacityConnector) Connect(id string) interfaces.Link {
	return NewPerfectLink(c, c.capacityMap.Connector(id))
}
