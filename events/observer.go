package events

import (
	"runtime"
)

// An observer is a special type of receiver that is used to monitor events
// related to a observer. When the simulation pops an event, all registered
// observers are notified. Registering an observer has priority over processing
// events, thus allowing registration at any moment.
type Observer interface {
	Receiver
}

// A passive observer keeps all the events that were received by it in a
// a channel. Moreover, it allows proxying the receival of an event by using
// using the Decorable interface. This means that the event can be processed
// before being put into the Recv channel.
type PassiveObserver interface {
	Decorable
	Observer

	Recv() <-chan interface{}
}

// An active observer is an observer which allows for a proxy. It's intent for
// usage is allowing for calling a Proxy at a higher level.
type ActiveObserver interface {
	Decorable
	Observer
}

// An active event observer calls the proxy on receival of an event.
type ActiveEventObserver struct {
	*Decorator
	receiver Receiver
}

func NewActiveEventObserver(receiver Receiver) *ActiveEventObserver {
	return &ActiveEventObserver{
		Decorator: NewDecorator(),
		receiver:  receiver,
	}
}

func (o *ActiveEventObserver) Receive(e *Event) *Event {
	if e.Receiver() == o.receiver {
		o.Proxy(e)
	}
	return nil
}

const maxGlobalObserverQueue int = 1000
const maxObserverQueue int = 1000

// A passive event observer is a PassiveObserver that calls a proxy function
// before delivering a message to the Recv channel.
type PassiveEventObserver struct {
	*Decorator

	receiver Receiver
	observer chan interface{}
}

func NewPassiveEventObserver(receiver Receiver) *PassiveEventObserver {
	return &PassiveEventObserver{
		Decorator: NewDecorator(),

		observer: make(chan interface{}, maxObserverQueue),
		receiver: receiver,
	}
}

func (o *PassiveEventObserver) Recv() <-chan interface{} {
	return o.observer
}

func (o *PassiveEventObserver) Receive(e *Event) *Event {
	if e.Receiver() == o.receiver {
		// Call the proxy on the received message before delivering it.
		// The Proxy can be set by upper layers via SetProxy(Proxy).
		deliver := o.Proxy(e)

		if deliver != nil {
			o.observer <- deliver
		}
	}
	return nil
}

// GlobalEventObserver is a passive observer which monitors all events. These
// events can be used for monitoring and logging purposes.
type GlobalEventObserver struct {
	*Decorator
	observer chan interface{}
}

func NewGlobalEventObserver() *GlobalEventObserver {
	return &GlobalEventObserver{
		Decorator: NewDecorator(),
		observer:  make(chan interface{}, maxObserverQueue),
	}
}

func (o *GlobalEventObserver) Recv() <-chan interface{} {
	return o.observer
}

func (o *GlobalEventObserver) Receive(e *Event) *Event {
	select {
	case o.observer <- o.Proxy(e):
	default:
		// If we can't register the observer, let someone else to run
		// to break the livelock.
		runtime.Gosched()
		o.Receive(e)
	}
	return nil
}
