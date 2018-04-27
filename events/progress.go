package events

import (
  . "github.com/danalex97/Speer/interfaces"
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
