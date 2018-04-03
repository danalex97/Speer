package capacity

import (
  "container/list"
  . "github.com/danalex97/Speer/overlay"
)

/**
 Assumptions:
   - equal link share
   - end-capacity is the bottleneck
 */

type Link interface {
  Upload(Data)
  Download() <-chan Data
}

const MaxConnections int = 50

type Data struct {
  id   string
  size int
}

type PerfectLink struct {
  bandwidth int
  from      UnreliableNode
  to        UnreliableNode
  queue     list.List
  time      int
  download  chan Data
}

func NewPerfectLink(from, to UnreliableNode, bandwidth int) Link {
  link := new(PerfectLink)

  link.bandwidth = bandwidth
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
