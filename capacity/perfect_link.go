package capacity

import (
  "container/list"
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

  From() Node
  To()   Node

  Up()   int
  Down() int
}

type Node interface {
}

const MaxConnections int = 100

type Data struct {
  Id   string
  Size int
}

type PerfectLink struct {
  up        int
  down      int
  from      Node
  to        Node
  queue     list.List
  download  chan Data
}

func NewPerfectLink(from, to Node, up, down int) Link {
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

func (p *PerfectLink) From() Node {
  return p.from
}

func (p *PerfectLink) To() Node {
  return p.to
}
