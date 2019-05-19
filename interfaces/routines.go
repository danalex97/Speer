package interfaces

type RoutineCapabilities interface {
	// Routine(interval int, routine func()) Routine
	// Callback(interval int, routine func()) Callback
}

// Wrapers for routines that run at fixed intervals. The inital call of the
// routine is right after the routine is built. The next ones are done after
// the set interval. When the interval is changed via SetInterval the interval
// is changed only after the first call schduled by the old scheduling interval.
type Routine interface {
	Interval() int
	SetInterval(interval int)
}

// Wrapers for callback that is called once after `wait` units of time.
type Callback interface {
}
