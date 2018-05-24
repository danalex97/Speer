package events

import (
  . "github.com/danalex97/Speer/interfaces"
  "sync"
)

type ProgressProperty struct {
  Progress

  interval int
}

func NewProgressProperty(progress Progress, interval int) *ProgressProperty {
  return &ProgressProperty{
    Progress : progress,
    interval : interval,
  }
}

func (p *ProgressProperty) Receive(event *Event) *Event {
  p.Advance()

  return NewEvent(
    event.Timestamp() + p.interval,
    nil,
    p,
  )
}

// A progress property used to wait in groups.
type WGProgress struct {
  cond *sync.Cond

  mp   map[string]bool
  size int
  val  int

  start int
}

func NewWGProgress(start int) GroupProgress {
  return &WGProgress{
    cond  : &sync.Cond{L: &sync.Mutex{}},

    mp    : make(map[string]bool),
    start : start,
  }
}

func (p *WGProgress) Add() {
  p.cond.L.Lock()
  defer p.cond.L.Unlock()

  p.size++
}

func (p *WGProgress) Progress(id string) {
  p.cond.L.Lock()
  defer p.cond.L.Unlock()

  p.mp[id] = true
  if len(p.mp) >= p.val  {
    p.cond.Broadcast()
  }
}

func (p *WGProgress) Advance() {
  p.cond.L.Lock()
  defer p.cond.L.Unlock()

  p.val = p.size

  if p.size >= p.start {
    p.mp = make(map[string]bool)

    p.cond.Wait()
  }
}
