package metrics

import (
  . "github.com/danalex97/Speer/events"
  "runtime"
  "fmt"
)

type Metrics struct {
  events <-chan *Event
}

func NewMetrics(o EventObserver) *Metrics {
  metrics := new(Metrics)
  metrics.events = o.EventChan()
  return metrics
}

func (m *Metrics) Run() {
  select {
  case  <-m.events:
    fmt.Println("metrics")
  default:
    runtime.Gosched()
  }
}
