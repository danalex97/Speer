package interfaces

type RoutineCapabilities interface {
	Routine(interval int, routine func()) Routine
	Callback(timeout int, routine func()) Callback
}

// Wrapers for routines that run at fixed intervals. The inital call of the
// routine is right after the routine is built. The next ones are done after
// the set interval. When the interval is changed via SetInterval the interval
// is changed only after the first call schduled by the old scheduling interval.
// Routines can be stopped.
type Routine interface {
	Interval() int
	SetInterval(interval int)

	Stop()
}

// Wrapers for callback that is called once after `wait` units of time.
// Callbacks can be stopped as well.
type Callback interface {
	Stop()
}
