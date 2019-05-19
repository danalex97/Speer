package interfaces

// Wrapers for routines that run at fixed intervals.
type Routine interface {
	NewRoutine(interval int, routine func())

	Interval() int
	SetInterval(interval int)
}

// Wrapers for callback that is called once after `wait` units of time.
type Callback interface {
	NewCallback(interval int, routine func())
}
