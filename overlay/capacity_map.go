package overlay

import (
	. "github.com/danalex97/Speer/capacity"

	"sync"
)

// Allow access the CapacityConnector
type CapacityMap interface {
	AddConnector(id string, connector CapacityConnector)
	Connector(id string) CapacityConnector
}

// A capacity map that uses a capacity scheduler for
// allowing data transfer between nodes. Note the capacity
// scheduler is already thread-safe.
type scheduledCapacityMap struct {
	*sync.RWMutex

	capacityMap       map[string]CapacityConnector
	capacityScheduler Scheduler
}

func NewScheduledCapacityMap(scheduleInterval int) CapacityMap {
	return &scheduledCapacityMap{
		// the rwmutex is used to synchronize the access to the
		// capacity connector map
		RWMutex: new(sync.RWMutex),

		// allows accessing capacity connectors and capacity scheduling
		capacityMap:       make(map[string]CapacityConnector),
		capacityScheduler: NewScheduler(scheduleInterval),
	}
}

func (c *scheduledCapacityMap) AddConnector(
	id string,
	capacityConnector CapacityConnector,
) {
	c.Lock()
	defer c.Unlock()

	c.capacityMap[id] = capacityConnector
}

func (c *scheduledCapacityMap) Connector(id string) CapacityConnector {
	c.RLock()
	defer c.RUnlock()

	return c.capacityMap[id]
}
