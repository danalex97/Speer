package capacity

import (
  . "github.com/danalex97/Speer/interfaces"
  "container/list"
  "sync"
)

/**
 Assumptions:
   - equal link share & max capacity
   - no latency effects
   - end-capacity is the bottleneck
 */

const MaxConnections int = 100

type PerfectLink struct {
  *sync.Mutex

  from      Node
  to        Node
  queue     *list.List
  download  chan Data
}

func NewPerfectLink(from, to Node) Link {
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

func (p *PerfectLink) From() Node {
  return p.from
}

func (p *PerfectLink) To() Node {
  return p.to
}

func (p *PerfectLink) Clear() {
  p.Lock()
  defer p.Unlock()

  p.queue = list.New()
}
