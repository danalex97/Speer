package interfaces

// A Simulation interface exposed when using the SDK.
type ISimulation interface {
	Run()
	Stop()
	Time() int
}
