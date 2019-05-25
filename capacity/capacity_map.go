package capacity

import (
	"github.com/danalex97/Speer/events"
	"github.com/danalex97/Speer/interfaces"
	"github.com/danalex97/Speer/underlay"

	"sync"
)

// Allow access the CapacityConnector
type CapacityMap interface {
	AddConnector(id string, connector CapacityConnector)
	Connector(id string) CapacityConnector

	RegisterLink(link interfaces.Link)
	Start(simulation *underlay.NetworkSimulation)

	// A simulation Receiver associated with the CapacityMap. Allows
	// registering observers that follow the map.
	Receiver() events.Receiver
}

// A capacity map that uses a capacity scheduler for
// allowing data transfer between nodes. Note the capacity
// scheduler is already thread-safe.
type ScheduledCapacityMap struct {
	*sync.RWMutex

	capacityMap       map[string]CapacityConnector
	capacityScheduler Scheduler
}

func NewScheduledCapacityMap(scheduleInterval int) CapacityMap {
	return &ScheduledCapacityMap{
		// the rwmutex is used to synchronize the access to the
		// capacity connector map
		RWMutex: new(sync.RWMutex),

		// allows accessing capacity connectors and capacity scheduling
		capacityMap:       make(map[string]CapacityConnector),
		capacityScheduler: NewScheduler(scheduleInterval),
	}
}

func (c *ScheduledCapacityMap) AddConnector(
	id string,
	capacityConnector CapacityConnector,
) {
	c.Lock()
	defer c.Unlock()

	c.capacityMap[id] = capacityConnector
}

func (c *ScheduledCapacityMap) Connector(id string) CapacityConnector {
	c.RLock()
	defer c.RUnlock()

	return c.capacityMap[id]
}

func (c *ScheduledCapacityMap) RegisterLink(link interfaces.Link) {
	c.capacityScheduler.RegisterLink(link)
}

func (c *ScheduledCapacityMap) Start(sim *underlay.NetworkSimulation) {
	sim.Push(events.NewEvent(0, nil, c.capacityScheduler))
}

func (c *ScheduledCapacityMap) Receiver() events.Receiver {
	return c.capacityScheduler
}
