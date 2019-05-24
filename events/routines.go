package events

import (
	"github.com/danalex97/Speer/interfaces"
)

type (
	RoutineReceiver interface {
		interfaces.Routine
		Receiver
	}

	CallbackReceiver interface {
		interfaces.Callback
		Receiver
	}

	routine struct {
		interval int
		routine  func()
		stopped  bool
	}

	callback struct {
		routine func()
		stopped bool
	}
)

func NewRoutine(interval int, exec func()) RoutineReceiver {
	return &routine{
		interval: interval,
		routine:  exec,
		stopped:  false,
	}
}

func (r *routine) Interval() int {
	return r.interval
}

func (r *routine) SetInterval(interval int) {
	r.interval = interval
}

func (r *routine) Receive(event *Event) *Event {
	if r.stopped {
		return nil
	}
	r.routine()
	return NewEvent(
		event.Timestamp()+r.interval,
		nil,
		r,
	)
}

func (r *routine) Stop() {
	r.stopped = true
}

func NewCallback(exec func()) CallbackReceiver {
	return &callback{
		routine: exec,
		stopped: false,
	}
}

func (c *callback) Receive(event *Event) *Event {
	if c.stopped {
		return nil
	}
	c.routine()
	return nil
}

func (c *callback) Stop() {
	c.stopped = true
}
