package events

import (
  . "github.com/danalex97/Speer/interfaces"
)

// We define a Progress as an interface by a progress function which should be
// called by a routine when progress is being made and an advance function that
// should block until a certain progress property has been made. The Progress
// runs inside as Receiver which periodically enqueues an event which checks
// for a progress property being made.
//
// By using this simple mechanism we can define both progress and safety
// properties. This allows the possibility of using concurrent functions without
// degrading the correctness of the simulation.
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
