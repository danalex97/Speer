package capacity

import (
  "container/list"
  . "github.com/danalex97/Speer/overlay"
)

/**
 Assumptions:
   - equal link share & max capacity
   - no latency effects
   - end-capacity is the bottleneck
 */

type Link interface {
  Upload(Data)
  Download() <-chan Data

  Up()   int
  Down() int
}

const MaxConnections int = 100

type Data struct {
  id   string
  size int
}

type PerfectLink struct {
  up        int
  down      int
  from      UnreliableNode
  to        UnreliableNode
  queue     list.List
  download  chan Data
}

func NewPerfectLink(from, to UnreliableNode, up, down int) Link {
  link := new(PerfectLink)

  link.up = up
  link.down = down

  link.from = from
  link.to = to
  link.queue = *list.New()
  link.download = make(chan Data, MaxConnections)

  return link
}

func (p *PerfectLink) Upload(data Data) {
  p.queue.PushBack(data)
}

func (p *PerfectLink) Download() <-chan Data {
  return p.download
}

func (p *PerfectLink) Up() int {
  return p.up
}

func (p *PerfectLink) Down() int {
  return p.down
}
