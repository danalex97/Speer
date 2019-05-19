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
	}

	callback struct {
		routine func()
	}
)

func NewRoutine(interval int, exec func()) RoutineReceiver {
	return &routine{
		interval: interval,
		routine:  exec,
	}
}

func (r *routine) Interval() int {
	return r.interval
}

func (r *routine) SetInterval(interval int) {
	r.interval = interval
}

func (r *routine) Receive(event *Event) *Event {
	r.routine()
	return NewEvent(
		event.Timestamp()+r.interval,
		nil,
		r,
	)
}

func NewCallback(exec func ()) CallbackReceiver {
	return &callback{
		routine: exec,
	}
}

func (c *callback) Receive(event *Event) *Event {
	c.routine()
	return nil
}
