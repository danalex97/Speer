package capacity

import (
  . "github.com/danalex97/Speer/events"
  "sync"
)

type Scheduler interface {
  Receiver
  RegisterLink(Link)
  Schedule()
}

type scheduler struct {
  interval int

  cntMutex *sync.RWMutex

  cnt      int
  linkCnt  map[Link]int
}

func NewScheduler(interval int) Scheduler {
  s := new(scheduler)

  s.interval = interval
  s.cntMutex = new(sync.RWMutex)
  s.cnt      = 0
  s.linkCnt  = make(map[Link]int)

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

func (s *scheduler) Schedule() {
  s.cntMutex.RLock()
  defer s.cntMutex.RUnlock()

}

func (s *scheduler) RegisterLink(l Link) {
  s.cntMutex.RLock()
  defer s.cntMutex.RUnlock()

  s.linkCnt[l] = s.cnt
}
