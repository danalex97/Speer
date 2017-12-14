package overlay

import (
	. "github.com/danalex97/Speer/capacity"
	"github.com/danalex97/Speer/interfaces"
)

// Allow connecting between two nodes via scheduler
type CapacityConnector interface {
	interfaces.DataTransport
}

type PerfectCapacityConnector struct {
	capacityMap CapacityMap

	down int
	up   int
}

func NewCapacityConnector(
	up, down int,
	capacityMap CapacityMap,
) CapacityConnector {
	return &PerfectCapacityConnector{
		capacityMap: capacityMap,

		down: down,
		up:   up,
	}
}

func (c *PerfectCapacityConnector) Up() int {
	return c.up
}

func (c *PerfectCapacityConnector) Down() int {
	return c.down
}

func (c *PerfectCapacityConnector) Connect(id string) interfaces.Link {
	return NewPerfectLink(c, c.capacityMap.Connector(id))
}
