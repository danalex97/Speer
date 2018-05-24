package events

import (
  . "github.com/danalex97/Speer/interfaces"
  "runtime"
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
const progressChanSize int = 1000

type WGProgress struct {
  *sync.Mutex

  notif chan string
  size  int
}

func NewWGProgress() GroupProgress {
  return &WGProgress{
    Mutex : new(sync.Mutex),

    notif : make(chan string, progressChanSize),
    size  : 0,
  }
}

func (p *WGProgress) Add() {
  p.Lock()
  defer p.Unlock()

  p.size++
}

func (p *WGProgress) compact(id string) {
  p.Lock()
  defer p.Unlock()

  if len(p.notif) < maxRegisterQueue {
    select {
      case p.notif <- id:
      default:
    }
    return
  }

  notif := make(chan string, progressChanSize)
  mp    := make(map[string]bool)

  // fmt.Println("Before ", len(p.notif))

  for {
    select {
    case id := <-p.notif:
      if _, ok := mp[id]; !ok {
        mp[id] = true
      }
    default:
      for id, _ := range mp {
        notif <- id
      }
      // fmt.Println("After ", len(notif))
      p.notif = notif
      p.notif <- id

      return
    }
  }
}

func (p *WGProgress) Progress(id string) {
  select {
  case p.notif <- id:
  default:
    p.compact(id)
  }
}

func (p *WGProgress) Advance() {
  p.Lock()
  defer p.Unlock()

  target  := p.size
  current := 0

  mp := make(map[string]bool)

  for current < target {
    select {
    case id := <-p.notif:
      if _, ok := mp[id]; !ok {
        mp[id] = true
        current++
      }
    default:
      p.Unlock()
      runtime.Gosched()
      p.Lock()
    }
  }

  p.notif = make(chan string, progressChanSize)
}
