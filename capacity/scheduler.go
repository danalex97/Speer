package capacity

import (
  . "github.com/danalex97/Speer/structs"
  . "github.com/danalex97/Speer/events"
  "sync"
  "math"
)

type Scheduler interface {
  Receiver
  RegisterLink(Link)
  Schedule()
}

type scheduler struct {
  interval int

  cntMutex *sync.RWMutex

  cnt        int
  linkStatus map[Link]*status
}

type status struct {
  active   bool
  data     float64
  capacity float64
}

func NewScheduler(interval int) Scheduler {
  s := new(scheduler)

  s.interval   = interval
  s.cntMutex   = new(sync.RWMutex)
  s.cnt        = 0
  s.linkStatus = make(map[Link]*status)

  return s
}

func (s *scheduler) Receive(event *Event) *Event {
  s.cntMutex.RLock()
  defer s.cntMutex.RUnlock()

  s.Schedule()
  s.cnt += 1

  return NewEvent(
    event.Timestamp() + s.interval,
    nil,
    s,
  )
}

func (s *scheduler) updData() {
  for link, status := range s.linkStatus {
    link := link.(*PerfectLink)

    if status.active {
      status.data += status.capacity * float64(s.interval)

      data := link.queue.Front().Value.(Data)
      for float64(data.Size) <= status.data {
        // dequeue data
        status.data -= float64(data.Size)
        link.download <- data

        // remove data from queue
        link.queue.Remove(link.queue.Front())
        if link.queue.Len() == 0 {
          // if the queue is empty we don't use the link
          status.active = false
          status.data   = 0
          break
        }
      }
    } else {
      // If we have pending requests, make the link active
      if link.queue.Len() > 0 {
        status.active = true
        status.data   = 0
      }
    }
  }
}

type elem struct {
  link Link
  seq  int
}

func (s *scheduler) updCapacity() {
  pq  := NewPriorityQueue()
  seq := make(map[Link]int)

  in  := make(map[Node]map[Link]bool)
  out := make(map[Node]map[Link]bool)

  up   := make(map[Node]float64)
  down := make(map[Node]float64)

  // build in, out map
  for link, status := range s.linkStatus {
    if status.active {
      if _, ok := in[link.To()]; !ok {
        in[link.To()] = make(map[Link]bool)
        down[link.To()] = float64(link.To().Down())
      }
      if _, ok := out[link.From()]; !ok {
        out[link.From()] = make(map[Link]bool)
        up[link.From()] = float64(link.From().Up())
      }

      in[link.To()][link] = true
      out[link.From()][link] = true
    }
  }

  capacity := func(link Link) float64 {
    outDeg := len(out[link.From()])
    inDeg  := len(in[link.To()])

    upCap := up[link.From()] / float64(outDeg)
    downCap := down[link.To()] / float64(inDeg)

    cap := math.Min(upCap, downCap)
    return cap
  }

  // build pq
  for link, status := range s.linkStatus {
    if status.active {
      cap := capacity(link)

      seq[link] = 0
      pq.Push(Float(cap), &elem{link, seq[link]})
    }
  }

  for pq.Len() > 0 {
    front := pq.Pop()

    top := front.Value.(*elem)
    cap := float64(front.Key.(Float))

    link := top.link
    if top.seq < seq[link] {
      // the element is stale, so we can continue
      continue
    }

    // update the link
    status := s.linkStatus[link]
    status.capacity = cap

    // remove the smallest capacity link
    delete(out[link.From()], link)
    delete(in[link.To()], link)

    // update new link capacities
    update := func(l Link) {
      cap := capacity(l)

      seq[l] = seq[l] + 1
      pq.Push(Float(cap), &elem{l, seq[l]})

      // update the download and upload capacities
      up[l.From()] = up[l.From()] - cap
      down[l.To()] = down[l.To()] - cap
    }

    for l, _ := range out[link.From()] {
      update(l)
    }

    for l, _ := range in[link.To()] {
      update(l)
    }
  }
}

func (s *scheduler) Schedule() {
  s.cntMutex.RLock()
  defer s.cntMutex.RUnlock()

  s.updData()
  s.updCapacity()
}

func (s *scheduler) RegisterLink(l Link) {
  s.cntMutex.RLock()
  defer s.cntMutex.RUnlock()

  s.linkStatus[l] = &status{false, 0, 0}
}
