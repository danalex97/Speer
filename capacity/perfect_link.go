package capacity

import (
	"container/list"
	. "github.com/danalex97/Speer/interfaces"
	"sync"
)

const MaxConnections int = 100

// Assumptions:
//   - equal link share & max capacity
//   - no latency effects
//   - end-capacity is the bottleneck
type PerfectLink struct {
	*sync.Mutex

	from     NodeCapacity
	to       NodeCapacity
	queue    *list.List
	download chan Data
}

func NewPerfectLink(from, to NodeCapacity) Link {
	return &PerfectLink{
		new(sync.Mutex),
		from,
		to,
		list.New(),
		make(chan Data, MaxConnections),
	}
}

func (p *PerfectLink) Upload(data Data) {
	p.Lock()
	defer p.Unlock()

	p.queue.PushBack(data)
}

func (p *PerfectLink) Download() <-chan Data {
	return p.download
}

func (p *PerfectLink) From() NodeCapacity {
	return p.from
}

func (p *PerfectLink) To() NodeCapacity {
	return p.to
}

func (p *PerfectLink) Clear() {
	p.Lock()
	defer p.Unlock()

	p.queue = list.New()
}
