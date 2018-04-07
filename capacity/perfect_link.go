package capacity

import (
  . "github.com/danalex97/Speer/interfaces"
  "container/list"
)

/**
 Assumptions:
   - equal link share & max capacity
   - no latency effects
   - end-capacity is the bottleneck
 */

const MaxConnections int = 100

type PerfectLink struct {
  from      Node
  to        Node
  queue     *list.List
  download  chan Data
}

func NewPerfectLink(from, to Node) Link {
  link := new(PerfectLink)

  link.from = from
  link.to = to
  link.queue = list.New()
  link.download = make(chan Data, MaxConnections)

  return link
}

func (p *PerfectLink) Upload(data Data) {
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
