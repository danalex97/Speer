package interfaces

type ISimulation interface {
  Run()
  Stop()
  Time() int
}
